import { Logger, Module } from '@nestjs/common';
import { DrtCommunicationModule } from '../../common/services/drt-communication/drt-communication.module';
import { PinataService } from './pinata.service';

@Module({
  providers: [Logger, PinataService],
  imports: [DrtCommunicationModule],
  exports: [PinataService],
})
export class IpfsModule {}
