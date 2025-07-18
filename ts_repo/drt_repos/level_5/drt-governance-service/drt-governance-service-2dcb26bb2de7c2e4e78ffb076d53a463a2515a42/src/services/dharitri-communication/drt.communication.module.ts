import { Module } from '@nestjs/common';
import { ApiConfigService } from 'src/helpers/api.config.service';
import { MXApiService } from './drt.api.service';
import { MXGatewayService } from './drt.gateway.service';
import { MXProxyService } from './drt.proxy.service';
import { DynamicModuleUtils } from 'src/utils/dynamic.module.utils';

@Module({
    imports: [DynamicModuleUtils.getCacheModule()],
    providers: [
        MXProxyService,
        MXApiService,
        MXGatewayService,
        ApiConfigService,
    ],
    exports: [
        MXProxyService,
        MXApiService,
        MXGatewayService,
        ApiConfigService,
    ],
})
export class MXCommunicationModule {}
