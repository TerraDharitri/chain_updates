import { forwardRef, Module } from "@nestjs/common";
import { TokenModule } from "src/endpoints/tokens/token.module";
import { MoaFarmActionRecognizerService } from "./moa.farm.action.recognizer.service";
import { MoaPairActionRecognizerService } from "./moa.pair.action.recognizer.service";
import { TransactionActionMoaRecognizerService } from "./transaction.action.moa.recognizer.service";
import { MoaWrapActionRecognizerService } from "./moa.wrap.action.recognizer.service";
import { MoaDistributionActionRecognizerService } from "./moa.distribution.action.recognizer.service";
import { TransactionActionModule } from "../../transaction.action.module";
import { MoaLockedAssetActionRecognizerService } from "./moa.locked.asset.action.recognizer.service";
import { ApiConfigModule } from "src/common/api-config/api.config.module";
import { MoaModule } from "src/endpoints/moa/moa.module";

@Module({
  imports: [
    forwardRef(() => TokenModule),
    forwardRef(() => TransactionActionModule),
    ApiConfigModule,
    MoaModule.forRoot(),
  ],
  providers: [
    TransactionActionMoaRecognizerService,
    MoaPairActionRecognizerService,
    MoaFarmActionRecognizerService,
    MoaWrapActionRecognizerService,
    MoaDistributionActionRecognizerService,
    MoaLockedAssetActionRecognizerService,
  ],
  exports: [TransactionActionMoaRecognizerService],
})
export class TransactionActionMoaRecognizerModule { }
