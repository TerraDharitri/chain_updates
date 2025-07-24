import { Injectable } from "@nestjs/common";
import { TransactionAction } from "../../entities/transaction.action";
import { TransactionActionCategory } from "../../entities/transaction.action.category";
import { TransactionMetadata } from "../../entities/transaction.metadata";
import { MoaFunction } from "./entities/moa.function.options";
import { MoaSettings } from "../../../../moa/entities/moa.settings";

@Injectable()
export class MoaDistributionActionRecognizerService {
  recognize(settings: MoaSettings, metadata: TransactionMetadata): TransactionAction | undefined {
    if (metadata.receiver === settings.distributionContract && metadata.functionName === MoaFunction.claimLockedAssets) {
      return this.getClaimLockedAssetsAction();
    }

    return undefined;
  }

  private getClaimLockedAssetsAction(): TransactionAction {
    const result = new TransactionAction();
    result.category = TransactionActionCategory.moa;
    result.name = MoaFunction.claimLockedAssets;
    result.description = 'Claim locked assets';

    return result;
  }
}
