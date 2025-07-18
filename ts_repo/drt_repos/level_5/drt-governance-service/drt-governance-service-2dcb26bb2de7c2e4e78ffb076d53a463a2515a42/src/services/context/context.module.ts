import { Module } from '@nestjs/common';
import { DRTCommunicationModule } from '../dharitri-communication/drt.communication.module';
import { ContextGetterService } from './context.getter.service';

@Module({
    imports: [DRTCommunicationModule],
    providers: [ContextGetterService],
    exports: [ContextGetterService],
})
export class ContextModule {}
