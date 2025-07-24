import { forwardRef, Global, Logger, Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config/dist';
import { ScheduleModule } from '@nestjs/schedule';
import { DrtCommunicationModule } from './common/services/drt-communication/drt-communication.module';
import { CacheModule } from './common/services/caching/caching.module';
import { ApiConfigModule } from './modules/common/api-config/api.config.module';

@Global()
@Module({
  imports: [forwardRef(() => CacheModule), ScheduleModule.forRoot(), ConfigModule, ApiConfigModule, DrtCommunicationModule],
  exports: [DrtCommunicationModule, CacheModule, Logger, ApiConfigModule],
  providers: [Logger],
})
export class CommonModule {}
