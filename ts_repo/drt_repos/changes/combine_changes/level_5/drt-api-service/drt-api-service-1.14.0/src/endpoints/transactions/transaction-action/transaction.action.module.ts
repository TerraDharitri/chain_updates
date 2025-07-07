import { forwardRef, Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";
import { TokenModule } from "src/endpoints/tokens/token.module";
import { TransactionActionDcdtNftRecognizerService } from "./recognizers/dcdt/transaction.action.dcdt.nft.recognizer.service";
import { TransactionActionService } from "./transaction.action.service";
import { TransactionActionMoaRecognizerModule } from "./recognizers/moa/transaction.action.moa.recognizer.module";
import { StakeActionRecognizerService } from "./recognizers/staking/transaction.action.stake.recognizer.service";
import { SCCallActionRecognizerService } from "./recognizers/sc-calls/transaction.action.sc-calls.recognizer.service";
import { ProviderModule } from "src/endpoints/providers/provider.module";
import { IdentitiesModule } from "src/endpoints/identities/identities.module";
import { MetabondingActionRecognizerService } from "./recognizers/moa/moa.metabonding.action.recognizer.service";

@Module({
  imports: [
    forwardRef(() => TokenModule),
    ConfigModule,
    TransactionActionMoaRecognizerModule,
    forwardRef(() => ProviderModule),
    forwardRef(() => IdentitiesModule),
  ],
  providers: [
    TransactionActionService,
    TransactionActionDcdtNftRecognizerService,
    StakeActionRecognizerService,
    SCCallActionRecognizerService,
    MetabondingActionRecognizerService,
  ],
  exports: [
    TransactionActionService,
    TransactionActionDcdtNftRecognizerService,
  ],
})
export class TransactionActionModule { }
