import { Module } from '@nestjs/common';
import { CommonAppModule } from 'src/common.app.module';
import { ContextModule } from 'src/services/context/context.module';
import { DRTCommunicationModule } from 'src/services/dharitri-communication/drt.communication.module';
import { TokenModule } from '../tokens/token.module';
import { EnergyAbiService } from './services/energy.abi.service';
import { EnergyComputeService } from './services/energy.compute.service';
import { EnergyService } from './services/energy.service';
import { EnergySetterService } from './services/energy.setter.service';

@Module({
    imports: [
        CommonAppModule,
        DRTCommunicationModule,
        ContextModule,
        TokenModule,
    ],
    providers: [
        EnergyService,
        EnergyAbiService,
        EnergySetterService,
        EnergyComputeService,
    ],
    exports: [
        EnergyAbiService,
        EnergySetterService,
        EnergyComputeService,
        EnergyService,
    ],
})
export class EnergyModule {}
