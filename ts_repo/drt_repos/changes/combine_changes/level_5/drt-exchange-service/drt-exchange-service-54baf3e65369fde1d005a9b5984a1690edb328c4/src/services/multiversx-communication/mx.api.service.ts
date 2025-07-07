import { constantsConfig, drtConfig } from '../../config';
import { Inject, Injectable } from '@nestjs/common';
import { DcdtToken } from 'src/modules/tokens/models/dcdtToken.model';
import { NftCollection } from 'src/modules/tokens/models/nftCollection.model';
import { NftToken } from 'src/modules/tokens/models/nftToken.model';
import Agent, { HttpsAgent } from 'agentkeepalive';
import { WINSTON_MODULE_PROVIDER } from 'nest-winston';
import { Logger } from 'winston';
import { PerformanceProfiler } from '../../utils/performance.profiler';
import { MetricsCollector } from '../../utils/metrics.collector';
import { Stats } from '../../models/stats.model';
import { ApiConfigService } from 'src/helpers/api.config.service';
import { ApiNetworkProvider } from '@terradharitri/sdk-network-providers/out';
import {
    isDcdtToken,
    isDcdtTokenValid,
    isNftCollection,
    isNftCollectionValid,
} from 'src/utils/token.type.compare';
import { PendingExecutor } from 'src/utils/pending.executor';
import { MXProxyService } from './drt.proxy.service';
import { ContextTracker } from '@terradharitri/sdk-nestjs-common';

type GenericGetArgs = {
    methodName: string;
    resourceUrl: string;
    retries?: number;
};

@Injectable()
export class MXApiService {
    private readonly apiProvider: ApiNetworkProvider;
    private genericGetExecutor: PendingExecutor<GenericGetArgs, any>;

    constructor(
        private readonly apiConfigService: ApiConfigService,
        @Inject(WINSTON_MODULE_PROVIDER) private readonly logger: Logger,
        private readonly drtProxy: MXProxyService,
    ) {
        const keepAliveOptions = {
            maxSockets: drtConfig.keepAliveMaxSockets,
            maxFreeSockets: drtConfig.keepAliveMaxFreeSockets,
            timeout: this.apiConfigService.getKeepAliveTimeoutDownstream(),
            freeSocketTimeout: drtConfig.keepAliveFreeSocketTimeout,
            keepAlive: true,
        };
        const httpAgent = new Agent(keepAliveOptions);
        const httpsAgent = new HttpsAgent(keepAliveOptions);

        this.apiProvider = new ApiNetworkProvider(
            this.apiConfigService.getApiUrl(),
            {
                timeout: drtConfig.proxyTimeout,
                httpAgent: drtConfig.keepAlive ? httpAgent : null,
                httpsAgent: drtConfig.keepAlive ? httpsAgent : null,
                headers: {
                    origin: 'xExchangeService',
                },
                clientName: 'xExchangeService',
            },
        );
        this.genericGetExecutor = new PendingExecutor(
            async (getGenericArgs: GenericGetArgs) =>
                await this.doGetGeneric(
                    getGenericArgs.methodName,
                    getGenericArgs.resourceUrl,
                    getGenericArgs.retries,
                ),
        );
    }

    getService(): ApiNetworkProvider {
        return this.apiProvider;
    }

    private delay(ms: number) {
        return new Promise((resolve) => setTimeout(resolve, ms));
    }

    async doGetGeneric<T>(
        name: string,
        resourceUrl: string,
        retries = 1,
    ): Promise<T> {
        const profiler = new PerformanceProfiler(`${name} ${resourceUrl}`);
        try {
            const context = ContextTracker.get();
            if (
                this.apiConfigService.isDeephistoryActive() &&
                context &&
                context.deepHistoryTimestamp
            ) {
                resourceUrl = resourceUrl.includes('?')
                    ? `${resourceUrl}&timestamp=${context.deepHistoryTimestamp}`
                    : `${resourceUrl}?timestamp=${context.deepHistoryTimestamp}`;
            }

            const response = await this.getService().doGetGeneric(resourceUrl);
            profiler.stop();
            MetricsCollector.setApiCall(
                name,
                'dex-service',
                200,
                profiler.duration,
            );
            return response;
        } catch (error) {
            if (
                error.inner.isAxiosError &&
                error.inner.code === 'ECONNABORTED' &&
                retries < 3
            ) {
                await this.delay(500 * retries);
                return await this.doGetGeneric(name, resourceUrl, retries + 1);
            }
            this.logger.error(`${error.message} after ${retries} retries`, {
                path: `${MXApiService.name}.${name}`,
                resourceUrl,
            });
            throw new Error(error);
        } finally {
            profiler.stop();

            MetricsCollector.setExternalCall(
                MXApiService.name,
                name,
                profiler.duration,
            );
        }
    }

    async getStats(): Promise<Stats> {
        const stats = await this.doGetGeneric<Stats>(
            this.getStats.name,
            'stats',
        );
        return new Stats(stats);
    }

    async getShardBlockCountInEpoch(
        epoch: number,
        shardId: number,
    ): Promise<number> {
        return await this.doGetGeneric<number>(
            this.getStats.name,
            `blocks/count?epoch=${epoch}&shard=${shardId}`,
        );
    }

    async getAccountStats(address: string): Promise<any | undefined> {
        return await this.doGetGeneric(
            this.getAccountStats.name,
            `accounts/${address}`,
        );
    }

    async getToken(tokenID: string): Promise<DcdtToken> {
        try {
            const rawToken = await this.doGetGeneric<DcdtToken>(
                this.getToken.name,
                `tokens/${tokenID}`,
            );
            const dcdtToken = new DcdtToken(rawToken);
            if (!isDcdtToken(dcdtToken)) {
                return undefined;
            }

            if (!isDcdtTokenValid(dcdtToken)) {
                const gatewayToken = await this.drtProxy
                    .getService()
                    .getDefinitionOfFungibleToken(tokenID);
                dcdtToken.identifier = gatewayToken.identifier;
                dcdtToken.decimals = gatewayToken.decimals;
            }

            return dcdtToken;
        } catch (error) {
            return undefined;
        }
    }

    async getNftCollection(tokenID: string): Promise<NftCollection> {
        try {
            const rawCollection = await this.doGetGeneric(
                this.getNftCollection.name,
                `collections/${tokenID}`,
            );
            const collection = new NftCollection(rawCollection);
            if (!isNftCollection(collection)) {
                return undefined;
            }
            if (!isNftCollectionValid(collection)) {
                const gatewayCollection = await this.drtProxy
                    .getService()
                    .getDefinitionOfTokenCollection(tokenID);
                collection.decimals = gatewayCollection.decimals;
            }
            return collection;
        } catch (error) {
            return undefined;
        }
    }

    async getTokensCountForUser(address: string): Promise<number> {
        return this.doGetGeneric<number>(
            this.getTokensCountForUser.name,
            `accounts/${address}/tokens/count`,
        );
    }

    async getNftsCountForUser(address: string): Promise<number> {
        return this.doGetGeneric<number>(
            this.getNftsCountForUser.name,
            `accounts/${address}/nfts/count`,
        );
    }

    async getTokensForUser(
        address: string,
        from = 0,
        size = 100,
    ): Promise<DcdtToken[]> {
        const profiler = new PerformanceProfiler('user_tokens');
        const userTokens = await this.doGetGeneric<DcdtToken[]>(
            this.getTokensForUser.name,
            `accounts/${address}/tokens?from=${from}&size=${size}`,
        );
        profiler.stop();
        MetricsCollector.setApiCall(
            this.getTokensForUser.name,
            'dex-service',
            200,
            profiler.duration,
        );

        for (const token of userTokens) {
            if (!isDcdtTokenValid(token)) {
                const gatewayToken = await this.drtProxy
                    .getService()
                    .getDefinitionOfFungibleToken(token.identifier);
                token.decimals = gatewayToken.decimals;
            }
        }

        return userTokens;
    }

    async getTokenForUser(
        address: string,
        tokenID: string,
    ): Promise<DcdtToken> {
        return this.doGetGeneric<DcdtToken>(
            this.getTokenForUser.name,
            `accounts/${address}/tokens/${tokenID}`,
        );
    }

    async getTokenBalanceForUser(
        address: string,
        tokenID: string,
    ): Promise<string> {
        try {
            const token = await this.getTokenForUser(address, tokenID);
            return token.balance;
        } catch (error) {
            return '0';
        }
    }

    async getNftsForUser(
        address: string,
        type = 'MetaDCDT',
    ): Promise<NftToken[]> {
        const nfts: NftToken[] = await this.genericGetExecutor.execute({
            methodName: this.getNftsForUser.name,
            resourceUrl: `accounts/${address}/nfts?type=${type}&size=${constantsConfig.MAX_USER_NFTS}&fields=identifier,collection,ticker,decimals,timestamp,attributes,nonce,type,name,creator,royalties,uris,url,tags,balance,assets`,
        });

        for (const nft of nfts) {
            if (!isNftCollectionValid(nft)) {
                const gatewayCollection = await this.drtProxy
                    .getService()
                    .getDefinitionOfTokenCollection(nft.collection);
                nft.decimals = gatewayCollection.decimals;
            }
        }

        return nfts;
    }

    async getNftsAttributesForUser(
        address: string,
        type = 'MetaDCDT',
        identifiers: string[],
    ): Promise<string[]> {
        if (identifiers.length === 0) {
            return [];
        }
        const nfts = await this.genericGetExecutor.execute({
            methodName: this.getNftsAttributesForUser.name,
            resourceUrl: `accounts/${address}/nfts?type=${type}&fields=attributes&identifiers=${identifiers.join(
                ',',
            )}`,
        });
        return nfts.map((nft) => nft.attributes);
    }

    async getNftByTokenIdentifier(
        address: string,
        nftIdentifier: string,
    ): Promise<NftToken> {
        return await this.doGetGeneric<NftToken>(
            this.getNftByTokenIdentifier.name,
            `accounts/${address}/nfts/${nftIdentifier}?fields=identifier,collection,ticker,decimals,timestamp,attributes,nonce,type,name,creator,royalties,uris,url,tags,balance,assets`,
        );
    }

    async getNftAttributesByTokenIdentifier(
        address: string,
        nftIdentifier: string,
    ): Promise<string> {
        const response = await this.doGetGeneric<NftToken>(
            this.getNftAttributesByTokenIdentifier.name,
            `accounts/${address}/nfts/${nftIdentifier}?fields=attributes`,
        );
        return response.attributes;
    }

    async getCurrentNonce(shardId: number): Promise<any> {
        return this.doGetGeneric(
            this.getCurrentNonce.name,
            `network/status/${shardId}`,
        );
    }

    async getCurrentBlockNonce(shardId: number): Promise<number> {
        const latestBlock = await this.doGetGeneric(
            this.getCurrentNonce.name,
            `blocks?size=1&shard=${shardId}`,
        );
        return latestBlock[0].nonce;
    }

    async getBlockByNonce(shardId: number, nonce: number): Promise<any> {
        const blocks = await this.doGetGeneric(
            this.getBlockByNonce.name,
            `blocks?nonce=${nonce}&shard=${shardId}`,
        );

        return blocks[0] ?? undefined;
    }

    async getShardTimestamp(shardId: number): Promise<number> {
        const latestShardBlock = await this.doGetGeneric(
            this.getShardTimestamp.name,
            `blocks?from=0&size=1&shard=${shardId}`,
        );
        return latestShardBlock[0].timestamp;
    }

    async getTransactions(
        after: number,
        before: number,
        receiverShard: number,
    ): Promise<any> {
        return await this.doGetGeneric(
            this.getTransactions.name,
            `transactions?receiverShard=${receiverShard}&after=${after}&before=${before}`,
        );
    }

    async getTransactionsWithOptions({
        sender,
        receiver,
        functionName,
    }): Promise<any> {
        return await this.doGetGeneric(
            this.getTransactions.name,
            `transactions?sender=${sender}&receiver=${receiver}&function=${functionName}`,
        );
    }
}
