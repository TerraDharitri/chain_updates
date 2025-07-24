import { Energy, EnergyType } from '@terradharitri/sdk-exchange';
import { Address, AddressValue, Interaction } from '@terradharitri/sdk-core';
import { Injectable } from '@nestjs/common';
import { DRTProxyService } from 'src/services/dharitri-communication/drt.proxy.service';
import { GenericAbiService } from 'src/services/generics/generic.abi.service';
import { GetOrSetCache } from 'src/helpers/decorators/caching.decorator';
import { Constants, ErrorLoggerAsync } from '@terradharitri/sdk-nestjs-common';
import { CacheTtlInfo } from 'src/services/caching/cache.ttl.info';
import { DRTApiService } from 'src/services/dharitri-communication/drt.api.service';
import { IEnergyAbiService } from './interfaces';

@Injectable()
export class EnergyAbiService
    extends GenericAbiService
    implements IEnergyAbiService
{
    constructor(
        protected readonly drtProxy: DRTProxyService,
        private readonly drtAPI: DRTApiService,
    ) {
        super(drtProxy);
    }

    @ErrorLoggerAsync()
    @GetOrSetCache({
        baseKey: 'energy',
        remoteTtl: Constants.oneMinute(),
    })
    async energyEntryForUser(userAddress: string): Promise<EnergyType> {
        return await this.getEnergyEntryForUserRaw(userAddress);
    }

    async getEnergyEntryForUserRaw(userAddress: string): Promise<EnergyType> {
        const contract = await this.drtProxy.getSimpleLockEnergySmartContract();
        const interaction: Interaction =
            contract.methodsExplicit.getEnergyEntryForUser([
                new AddressValue(Address.fromString(userAddress)),
            ]);

        const response = await this.getGenericData(interaction);
        const rawEnergy = response.firstValue.valueOf();
        return Energy.fromDecodedAttributes(rawEnergy).toJSON();
    }

    @ErrorLoggerAsync()
    @GetOrSetCache({
        baseKey: 'energy',
        remoteTtl: Constants.oneSecond(),
        localTtl: Constants.oneSecond(),
    })
    async energyAmountForUser(userAddress: string): Promise<string> {
        return await this.getEnergyAmountForUserRaw(userAddress);
    }

    async getEnergyAmountForUserRaw(userAddress: string): Promise<string> {
        const contract = await this.drtProxy.getSimpleLockEnergySmartContract();
        const interaction: Interaction =
            contract.methodsExplicit.getEnergyAmountForUser([
                new AddressValue(Address.fromString(userAddress)),
            ]);

        const response = await this.getGenericData(interaction);
        return response.firstValue.valueOf().toFixed();
    }

    @ErrorLoggerAsync()
    @GetOrSetCache({
        baseKey: 'energy',
        remoteTtl: CacheTtlInfo.ContractState.remoteTtl,
        localTtl: CacheTtlInfo.ContractState.localTtl,
    })
    async isPaused(): Promise<boolean> {
        return await this.isPausedRaw();
    }

    async isPausedRaw(): Promise<boolean> {
        const contract = await this.drtProxy.getSimpleLockEnergySmartContract();
        const interaction = contract.methodsExplicit.isPaused();
        const response = await this.getGenericData(interaction);

        return response.firstValue.valueOf();
    }
}
