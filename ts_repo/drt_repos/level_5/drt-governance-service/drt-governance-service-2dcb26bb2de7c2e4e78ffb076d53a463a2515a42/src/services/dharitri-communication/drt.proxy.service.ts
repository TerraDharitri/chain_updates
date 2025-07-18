import { AbiRegistry, Address, SmartContract } from '@terradharitri/sdk-core';
import { Inject, Injectable } from '@nestjs/common';
import { abiConfig, drtConfig, scAddress } from '../../config';
import Agent, { HttpsAgent } from 'agentkeepalive';
import { WINSTON_MODULE_PROVIDER } from 'nest-winston';
import { Logger } from 'winston';
import { ProxyNetworkProviderProfiler } from '../../helpers/proxy.network.provider.profiler';
import { ApiConfigService } from 'src/helpers/api.config.service';
import { promises } from 'fs';
import { GovernanceType } from '../../utils/governance';

@Injectable()
export class DRTProxyService {
    private readonly proxy: ProxyNetworkProviderProfiler;
    private static smartContracts: SmartContract[];

    constructor(
        private readonly apiConfigService: ApiConfigService,
        @Inject(WINSTON_MODULE_PROVIDER) private readonly logger: Logger,
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

        this.proxy = new ProxyNetworkProviderProfiler(
            this.apiConfigService.getApiUrl(),
            {
                timeout: drtConfig.proxyTimeout,
                httpAgent: drtConfig.keepAlive ? httpAgent : null,
                httpsAgent: drtConfig.keepAlive ? httpsAgent : null,
                headers: {
                    origin: 'MaiarExchangeService',
                },
            },
        );

        DRTProxyService.smartContracts = [];
    }

    getService(): ProxyNetworkProviderProfiler {
        return this.proxy;
    }

    async getAddressShardID(address: string): Promise<number> {
        const response = await this.getService().doGetGeneric(
            `address/${address}/shard`,
        );
        return response.shardID;
    }
    async getLockedAssetFactorySmartContract(): Promise<SmartContract> {
        return this.getSmartContract(
            scAddress.lockedAssetAddress,
            abiConfig.lockedAssetFactory,
            'LockedAssetFactory',
        );
    }

    async getSimpleLockEnergySmartContract(): Promise<SmartContract> {
        return this.getSmartContract(
            scAddress.simpleLockEnergy,
            abiConfig.simpleLockEnergy,
            'SimpleLockEnergy',
        );
    }

    async getGovernanceSmartContract(
        governanceAddress: string,
        type: GovernanceType,
    ): Promise<SmartContract> {
        return this.getSmartContract(
            governanceAddress,
            abiConfig.governance[type],
            'GovernanceV2',
        );
    }

    async getSmartContract(
        contractAddress: string,
        contractAbiPath: string,
        contractInterface: string,
    ): Promise<SmartContract> {
        const key = `${contractInterface}.${contractAddress}`;
        return (
            DRTProxyService.smartContracts[key] ||
            this.createSmartContract(
                contractAddress,
                contractAbiPath,
                contractInterface,
            )
        );
    }

    private async createSmartContract(
        contractAddress: string,
        contractAbiPath: string,
        contractInterface: string,
    ): Promise<SmartContract> {
        const jsonContent: string = await promises.readFile(contractAbiPath, {
            encoding: 'utf8',
        });
        const json = JSON.parse(jsonContent);
        const newSC = new SmartContract({
            address: Address.fromString(contractAddress),
            abi: AbiRegistry.create(json),
        });
        const key = `${contractInterface}.${contractAddress}`;
        DRTProxyService.smartContracts[key] = newSC;
        return newSC;
    }
}
