import { ContractQueryResponse } from '@terradharitri/sdk-network-providers/out';
import {
    Interaction,
    Query,
    ResultsParser,
    TypedOutcomeBundle,
} from '@terradharitri/sdk-core';
import { PendingExecutor } from 'src/utils/pending.executor';
import { DRTProxyService } from '../dharitri-communication/drt.proxy.service';

export class GenericAbiService {
    private queryExecutor: PendingExecutor<Query, ContractQueryResponse>;

    constructor(protected readonly drtProxy: DRTProxyService) {
        this.queryExecutor = new PendingExecutor(
            async (query: Query) =>
                await this.drtProxy.getService().queryContract(query),
        );
    }

    async getGenericData(
        interaction: Interaction,
    ): Promise<TypedOutcomeBundle> {
        const query = interaction.check().buildQuery();
        const queryResponse = await this.queryExecutor.execute(query);
        const endpointDefinition = interaction.getEndpoint();
        return new ResultsParser().parseQueryResponse(
            queryResponse,
            endpointDefinition,
        );
    }
}
