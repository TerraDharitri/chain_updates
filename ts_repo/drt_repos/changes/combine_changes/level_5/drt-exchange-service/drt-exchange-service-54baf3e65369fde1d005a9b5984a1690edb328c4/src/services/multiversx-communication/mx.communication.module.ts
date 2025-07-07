import { Module } from '@nestjs/common';
import { ApiConfigService } from 'src/helpers/api.config.service';
import { MXApiService } from './drt.api.service';
import { MXDataApiService } from './drt.data.api.service';
import { MXGatewayService } from './drt.gateway.service';
import { MXProxyService } from './drt.proxy.service';
import { DynamicModuleUtils } from 'src/utils/dynamic.module.utils';

@Module({
    imports: [DynamicModuleUtils.getCacheModule()],
    providers: [
        MXProxyService,
        MXApiService,
        MXGatewayService,
        MXDataApiService,
        ApiConfigService,
    ],
    exports: [
        MXProxyService,
        MXApiService,
        MXGatewayService,
        MXDataApiService,
        ApiConfigService,
    ],
})
export class MXCommunicationModule {}
