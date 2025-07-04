import { Injectable } from "@nestjs/common";
import { TransactionAction } from "../../entities/transaction.action";
import { TransactionActionCategory } from "../../entities/transaction.action.category";
import { TransactionMetadata } from "../../entities/transaction.metadata";
import { TransactionActionDcdtNftRecognizerService } from "../dcdt/transaction.action.dcdt.nft.recognizer.service";
import { MoaFunction } from "./entities/moa.function.options";
import { MoaSettings } from "../../../../moa/entities/moa.settings";

@Injectable()
export class MoaFarmActionRecognizerService {
  constructor(
    private readonly transactionActionDcdtNftRecognizerService: TransactionActionDcdtNftRecognizerService,
  ) { }

  recognize(settings: MoaSettings, metadata: TransactionMetadata): TransactionAction | undefined {
    if (!settings.farmContracts.includes(metadata.receiver)) {
      return undefined;
    }

    switch (metadata.functionName) {
      case MoaFunction.enterFarm:
      case MoaFunction.enterFarmProxy:
        return this.getFarmAction(metadata, MoaFunction.enterFarm, 'Enter farm with');
      case MoaFunction.enterFarmAndLockRewards:
      case MoaFunction.enterFarmAndLockRewardsProxy:
        return this.getFarmAction(metadata, MoaFunction.enterFarm, 'Enter farm and lock rewards with');
      case MoaFunction.exitFarm:
      case MoaFunction.exitFarmProxy:
        return this.getFarmAction(metadata, MoaFunction.exitFarm, 'Exit farm with');
      case MoaFunction.claimRewards:
      case MoaFunction.claimRewardsProxy:
        return this.getFarmAction(metadata, MoaFunction.claimRewards, 'Claim rewards for');
      case MoaFunction.compoundRewards:
      case MoaFunction.compoundRewardsProxy:
        return this.getFarmAction(metadata, MoaFunction.compoundRewards, 'Reinvest rewards for');
      case MoaFunction.stakeFarm:
      case MoaFunction.stakeFarmProxy:
        return this.getFarmAction(metadata, MoaFunction.enterFarm, 'Stake farm with');
      case MoaFunction.stakeFarmTokens:
      case MoaFunction.stakeFarmTokensProxy:
        return this.getFarmAction(metadata, MoaFunction.enterFarm, 'Stake farm tokens with');
      case MoaFunction.unstakeFarm:
      case MoaFunction.unstakeFarmProxy:
        return this.getFarmAction(metadata, MoaFunction.exitFarm, 'Unstake farm with');
      case MoaFunction.unstakeFarmTokens:
      case MoaFunction.unstakeFarmTokensProxy:
        return this.getFarmAction(metadata, MoaFunction.exitFarm, 'Unstake farm tokens with');
      case MoaFunction.claimDualYield:
      case MoaFunction.claimDualYieldProxy:
        return this.getFarmAction(metadata, MoaFunction.claimRewards, 'Claim dual yield for');
      case MoaFunction.unbondFarm:
        return this.getFarmAction(metadata, MoaFunction.unbondFarm, 'Unbond farm with');
      default:
        return undefined;
    }
  }

  private getFarmAction(metadata: TransactionMetadata, name: string, action: string): TransactionAction | undefined {
    return this.transactionActionDcdtNftRecognizerService.getMultiTransferActionWithTicker(metadata, TransactionActionCategory.moa, name, action);
  }
}
