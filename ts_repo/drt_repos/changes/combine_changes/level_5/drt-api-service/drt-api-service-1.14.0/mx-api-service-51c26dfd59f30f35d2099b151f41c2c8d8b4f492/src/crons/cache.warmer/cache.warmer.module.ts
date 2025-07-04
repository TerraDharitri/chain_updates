import { Module } from '@nestjs/common';
import { ScheduleModule } from '@nestjs/schedule';
import { CacheWarmerService } from './cache.warmer.service';
import { EndpointsServicesModule } from '../../endpoints/endpoints.services.module';
import { KeybaseModule } from 'src/common/keybase/keybase.module';
import { MoaModule } from 'src/endpoints/moa/moa.module';
import { AssetsModule } from 'src/common/assets/assets.module';
import { DynamicModuleUtils } from 'src/utils/dynamic.module.utils';
import { NftCronModule } from '../nft/nft.cron.module';
import { GuestCacheWarmer } from '@terradharitri/sdk-nestjs-cache';
import { PluginModule } from 'src/plugins/plugin.module';
import { TpsWarmerService } from '../tps/tps-warmer.service';
import { TpsModule } from 'src/endpoints/tps/tps.module';

@Module({
  imports: [
    ScheduleModule.forRoot(),
    EndpointsServicesModule,
    KeybaseModule,
    MoaModule.forRoot(),
    DynamicModuleUtils.getRedisCacheModule(),
    AssetsModule,
    NftCronModule,
    PluginModule,
    TpsModule,
  ],
  providers: [
    DynamicModuleUtils.getPubSubService(),
    GuestCacheWarmer,
    CacheWarmerService,
    TpsWarmerService,
  ],
})
export class CacheWarmerModule { }
