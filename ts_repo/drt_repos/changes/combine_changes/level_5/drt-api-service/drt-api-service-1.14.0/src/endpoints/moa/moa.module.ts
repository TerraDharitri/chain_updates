import { DynamicModule, Module, Provider, Type } from "@nestjs/common";
import configuration from "config/configuration";
import { GraphQlModule } from "src/common/graphql/graphql.module";
import { DynamicModuleUtils } from "src/utils/dynamic.module.utils";
import { MoaEconomicsService } from "./moa.economics.service";
import { MoaFarmService } from "./moa.farm.service";
import { MoaPairService } from "./moa.pair.service";
import { MoaSettingsService } from "./moa.settings.service";
import { MoaTokenChartsService } from "./moa.token.charts.service";
import { MoaTokenService } from "./moa.token.service";
import { MoaWarmerService } from "./moa.warmer.service";

@Module({})
export class MoaModule {
  static forRoot(): DynamicModule {
    const providers: (Type<any> | Provider<any>)[] = [
      DynamicModuleUtils.getPubSubService(),
      MoaEconomicsService,
      MoaSettingsService,
      MoaPairService,
      MoaTokenService,
      MoaFarmService,
      MoaTokenChartsService,
    ];

    const isExchangeEnabled = configuration().features?.exchange?.enabled ?? false;
    if (isExchangeEnabled) {
      providers.push(MoaWarmerService);
    }

    return {
      module: MoaModule,
      imports: [
        GraphQlModule,
      ],
      providers,
      exports: [
        MoaEconomicsService,
        MoaPairService,
        MoaSettingsService,
        MoaTokenService,
        MoaFarmService,
        MoaTokenChartsService,
      ],
    };
  }
}
