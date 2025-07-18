import { Injectable } from '@nestjs/common';
import { Interaction } from '@terradharitri/sdk-core/out/smartcontracts/interaction';
import { DRTProxyService } from 'src/services/dharitri-communication/drt.proxy.service';
import { GenericAbiService } from 'src/services/generics/generic.abi.service';

@Injectable()
export class AbiLockedAssetService extends GenericAbiService {
    constructor(protected readonly drtProxy: DRTProxyService) {
        super(drtProxy);
    }

    async getLockedTokenID(): Promise<string> {
        const contract =
            await this.drtProxy.getLockedAssetFactorySmartContract();
        const interaction: Interaction =
            contract.methodsExplicit.getLockedAssetTokenId();
        const response = await this.getGenericData(interaction);
        return response.firstValue.valueOf().toString();
    }
}
