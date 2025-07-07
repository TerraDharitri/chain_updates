import {
    ContractQueryResponse,
    ProxyNetworkProvider,
} from '@terradharitri/sdk-network-providers';
import { PerformanceProfiler } from '../utils/performance.profiler';
import { MetricsCollector } from '../utils/metrics.collector';
import { IContractQuery } from '@terradharitri/sdk-network-providers/out/interface';
import { ContextTracker } from '@terradharitri/sdk-nestjs-common';
import { ApiConfigService } from './api.config.service';
import { NetworkProviderConfig } from '@terradharitri/sdk-core/out/networkProviders/networkProviderConfig';

export class ProxyNetworkProviderProfiler extends ProxyNetworkProvider {
    constructor(
        private readonly apiConfigService: ApiConfigService,
        url: string,
        config?: NetworkProviderConfig,
    ) {
        super(url, config);
    }

    async queryContract(query: IContractQuery): Promise<ContractQueryResponse> {
        const profiler = new PerformanceProfiler();

        const result = await super.queryContract(query);

        profiler.stop();

        MetricsCollector.setExternalCall(
            'vm.query',
            query.func.toString(),
            profiler.duration,
        );

        return result;
    }

    async doPostGeneric(resourceUrl: string, payload: any): Promise<any> {
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
        const response = await super.doPostGeneric(resourceUrl, payload);
        return response;
    }
}
