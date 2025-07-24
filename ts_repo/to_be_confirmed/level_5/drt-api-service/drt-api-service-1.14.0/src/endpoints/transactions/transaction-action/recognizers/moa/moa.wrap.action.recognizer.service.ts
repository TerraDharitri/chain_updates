import { Injectable } from "@nestjs/common";
import { DcdtType } from "src/endpoints/dcdt/entities/dcdt.type";
import { TransactionAction } from "../../entities/transaction.action";
import { TransactionActionCategory } from "../../entities/transaction.action.category";
import { TransactionMetadata } from "../../entities/transaction.metadata";
import { TransactionActionDcdtNftRecognizerService } from "../dcdt/transaction.action.dcdt.nft.recognizer.service";
import { MoaFunction } from "./entities/moa.function.options";
import { MoaSettings } from "../../../../moa/entities/moa.settings";
import { MoaSettingsService } from "../../../../moa/moa.settings.service";
import { NumberUtils } from "@terradharitri/sdk-nestjs-common";

@Injectable()
export class MoaWrapActionRecognizerService {
  constructor(
    private readonly transactionActionDcdtNftRecognizerService: TransactionActionDcdtNftRecognizerService,
    private readonly moaSettingsService: MoaSettingsService,
  ) { }

  recognize(settings: MoaSettings, metadata: TransactionMetadata): TransactionAction | undefined {
    if (!settings.wrapContracts.includes(metadata.receiver)) {
      return undefined;
    }

    switch (metadata.functionName) {
      case MoaFunction.wrapRewa:
        return this.getWrapAction(metadata);
      case MoaFunction.unwrapRewa:
        return this.getUnwrapAction(metadata);
      default:
        return undefined;
    }
  }

  private getWrapAction(metadata: TransactionMetadata): TransactionAction | undefined {
    const wrewaId = this.moaSettingsService.getWrewaId();
    if (!wrewaId) {
      return undefined;
    }

    const valueDenominated = NumberUtils.toDenominatedString(metadata.value);


    const result = new TransactionAction();
    result.category = TransactionActionCategory.moa;
    result.name = MoaFunction.wrapRewa;
    result.description = `Wrap ${valueDenominated} REWA`;
    result.arguments = {
      token: {
        type: DcdtType.FungibleDCDT,
        name: 'WrappedREWA',
        token: wrewaId,
        ticker: wrewaId.split('-')[0],
        decimals: 18,
        value: metadata.value.toString(),
      },
      receiver: metadata.receiver,
    };

    return result;
  }

  private getUnwrapAction(metadata: TransactionMetadata): TransactionAction | undefined {
    return this.transactionActionDcdtNftRecognizerService.getMultiTransferActionWithTicker(metadata, TransactionActionCategory.moa, MoaFunction.unwrapRewa, 'Unwrap');
  }
}
