import { drtConfig } from '../../config';
import { Inject, Injectable } from '@nestjs/common';
import { DcdtToken } from 'src/modules/tokens/models/dcdtToken.model';
import Agent, { HttpsAgent } from 'agentkeepalive';
import { WINSTON_MODULE_PROVIDER } from 'nest-winston';
import { Logger } from 'winston';
import { PerformanceProfiler } from '../../utils/performance.profiler';
import { MetricsCollector } from '../../utils/metrics.collector';
import { Stats } from '../../models/stats.model';
import { ApiConfigService } from 'src/helpers/api.config.service';
import { ApiNetworkProvider } from '@terradharitri/sdk-network-providers/out';
import { isDcdtToken, isDcdtTokenValid } from 'src/utils/token.type.compare';
import { PendingExecutor } from 'src/utils/pending.executor';
import { DRTProxyService } from './drt.proxy.service';

type GenericGetArgs = {
    methodName: string;
    resourceUrl: string;
    retries?: number;
};

@Injectable()
export class DRTApiService {
    private readonly apiProvider: ApiNetworkProvider;
    private genericGetExecutor: PendingExecutor<GenericGetArgs, any>;

    constructor(
        private readonly apiConfigService: ApiConfigService,
        @Inject(WINSTON_MODULE_PROVIDER) private readonly logger: Logger,
        private readonly drtProxy: DRTProxyService,
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
                    origin: 'GovernanceService',
                },
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
            return await this.getService().doGetGeneric(resourceUrl);
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
                path: `${DRTApiService.name}.${name}`,
                resourceUrl,
            });
            throw new Error(error);
        } finally {
            profiler.stop();

            MetricsCollector.setExternalCall(
                DRTApiService.name,
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

    async getTokenForUser(
        address: string,
        tokenID: string,
    ): Promise<DcdtToken> {
        return this.doGetGeneric<DcdtToken>(
            this.getTokenForUser.name,
            `accounts/${address}/tokens/${tokenID}`,
        );
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
}
