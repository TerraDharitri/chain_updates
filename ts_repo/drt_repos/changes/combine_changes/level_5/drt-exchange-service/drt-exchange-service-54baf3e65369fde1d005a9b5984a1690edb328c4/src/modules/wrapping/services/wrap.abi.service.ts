import { Injectable } from '@nestjs/common';
import { Interaction } from '@terradharitri/sdk-core/out/smartcontracts/interaction';
import { MXProxyService } from '../../../services/dharitri-communication/drt.proxy.service';
import { GenericAbiService } from 'src/services/generics/generic.abi.service';
import { GetOrSetCache } from 'src/helpers/decorators/caching.decorator';
import { CacheTtlInfo } from 'src/services/caching/cache.ttl.info';
import { ErrorLoggerAsync } from '@terradharitri/sdk-nestjs-common';

@Injectable()
export class WrapAbiService extends GenericAbiService {
    constructor(protected readonly drtProxy: MXProxyService) {
        super(drtProxy);
    }

    @ErrorLoggerAsync()
    @GetOrSetCache({
        baseKey: 'wrap',
        remoteTtl: CacheTtlInfo.TokenID.remoteTtl,
        localTtl: CacheTtlInfo.TokenID.localTtl,
    })
    async wrappedRewaTokenID(): Promise<string> {
        return this.getWrappedRewaTokenIDRaw();
    }

    async getWrappedRewaTokenIDRaw(): Promise<string> {
        const contract = await this.drtProxy.getWrapSmartContract();
        const interaction: Interaction =
            contract.methodsExplicit.getWrappedRewaTokenId();
        const response = await this.getGenericData(interaction);
        return response.firstValue.valueOf().toString();
    }
}
