import { Injectable } from "@nestjs/common";
import { TransactionAction } from "../../entities/transaction.action";
import { TransactionActionCategory } from "../../entities/transaction.action.category";
import { TransactionMetadata } from "../../entities/transaction.metadata";
import { TransactionActionDcdtNftRecognizerService } from "../dcdt/transaction.action.dcdt.nft.recognizer.service";
import { MoaFunction } from "./entities/moa.function.options";
import { MoaSettings } from "../../../../moa/entities/moa.settings";
import { MoaSettingsService } from "../../../../moa/moa.settings.service";
import { NumberUtils } from "@terradharitri/sdk-nestjs-common";

@Injectable()
export class MoaLockedAssetActionRecognizerService {
  constructor(
    private readonly moaSettingsService: MoaSettingsService,
    private readonly transactionActionDcdtNftRecognizerService: TransactionActionDcdtNftRecognizerService,
  ) { }

  recognize(settings: MoaSettings, metadata: TransactionMetadata): TransactionAction | undefined {
    if (metadata.receiver !== settings.lockedAssetContract) {
      return undefined;
    }

    switch (metadata.functionName) {
      case MoaFunction.lockAssets:
        return this.getAssetsAction(metadata, 'Lock');
      case MoaFunction.unlockAssets:
        const action = this.getAssetsAction(metadata, 'Unlock');
        if (action) {
          action.description = 'Unlock assets';
        }

        return action;
      case MoaFunction.mergeLockedAssetTokens:
        return this.getMergeLockedAssetTokens(metadata);
      default:
        return undefined;
    }
  }

  private getMergeLockedAssetTokens(metadata: TransactionMetadata): TransactionAction | undefined {
    const transfers = this.moaSettingsService.getTransfers(metadata);
    if (!transfers) {
      return undefined;
    }

    const value = transfers.sumBigInt(x => BigInt(x.value.toString()));
    const valueDenominated = NumberUtils.toDenominatedString(value);

    const description = `Merge ${transfers.length} LKMOA positions into a single LKMOA position of value ${valueDenominated}`;

    return this.transactionActionDcdtNftRecognizerService.getMultiTransferAction(metadata, TransactionActionCategory.moa, MoaFunction.mergeLockedAssetTokens, description);
  }

  private getAssetsAction(metadata: TransactionMetadata, action: string): TransactionAction | undefined {
    return this.transactionActionDcdtNftRecognizerService.getMultiTransferActionWithTicker(metadata, TransactionActionCategory.moa, metadata.functionName ?? '', action);
  }
}
