import { forwardRef, Inject, Injectable } from "@nestjs/common";
import { TransactionAction } from "../../entities/transaction.action";
import { TransactionActionCategory } from "../../entities/transaction.action.category";
import { TransactionMetadata } from "../../entities/transaction.metadata";
import { MoaFunction } from "./entities/moa.function.options";
import { MoaSettings } from "../../../../moa/entities/moa.settings";
import { TokenTransferService } from "src/endpoints/tokens/token.transfer.service";
import { MoaSettingsService } from "../../../../moa/moa.settings.service";
import { TransactionActionDcdtNftRecognizerService } from "../dcdt/transaction.action.dcdt.nft.recognizer.service";
import { BinaryUtils, NumberUtils } from "@terradharitri/sdk-nestjs-common";

@Injectable()
export class MoaPairActionRecognizerService {
  constructor(
    private readonly moaSettingsService: MoaSettingsService,
    @Inject(forwardRef(() => TokenTransferService))
    private readonly tokenTransferService: TokenTransferService,
    private readonly transactionActionDcdtNftRecognizerService: TransactionActionDcdtNftRecognizerService,
  ) { }

  async recognize(settings: MoaSettings, metadata: TransactionMetadata): Promise<TransactionAction | undefined> {
    if (!settings.pairContracts.includes(metadata.receiver) && settings.routerFactoryContract !== metadata.receiver) {
      return undefined;
    }

    switch (metadata.functionName) {
      case MoaFunction.swapTokensFixedInput:
      case MoaFunction.swapTokensFixedOutput:
        return await this.getSwapAction(metadata);
      case MoaFunction.addLiquidity:
      case MoaFunction.addLiquidityProxy:
        return this.getAddLiquidityAction(metadata);
      case MoaFunction.removeLiquidity:
      case MoaFunction.removeLiquidityProxy:
        return this.getRemoveLiquidityAction(metadata);
      case MoaFunction.multiPairSwap:
        return this.getMultiSwapAction(metadata);
      default:
        return undefined;
    }
  }

  private async getSwapAction(metadata: TransactionMetadata): Promise<TransactionAction | undefined> {
    const transfers = this.moaSettingsService.getTransfers(metadata);
    if (!transfers) {
      return undefined;
    }

    const pair1Properties = transfers[0].properties;
    if (!pair1Properties) {
      return undefined;
    }

    const value = transfers[0].value;
    const valueDenominated = NumberUtils.toDenominatedString(value, pair1Properties.decimals);

    const destinationTokenIdentifier = BinaryUtils.hexToString(metadata.functionArgs[0]);
    const destinationValue = BinaryUtils.hexToBigInt(metadata.functionArgs[1]);

    const pair2Properties = await this.tokenTransferService.getTokenTransferProperties({ identifier: destinationTokenIdentifier });
    if (!pair2Properties) {
      return undefined;
    }

    const destinationValueDenominated = NumberUtils.toDenominatedString(destinationValue, pair2Properties.decimals);

    metadata.transfers?.push({
      value: destinationValue,
      properties: pair2Properties,
    });

    let description = `Swap ${valueDenominated} ${pair1Properties.ticker} for a minimum of ${destinationValueDenominated} ${pair2Properties.ticker}`;
    if (metadata.functionName === MoaFunction.swapTokensFixedOutput) {
      description = `Swap a maximum of ${valueDenominated} ${pair1Properties.ticker} for ${destinationValueDenominated} ${pair2Properties.ticker}`;
    }

    return this.transactionActionDcdtNftRecognizerService.getMultiTransferAction(metadata, TransactionActionCategory.moa, 'swap', description);
  }

  private async getMultiSwapAction(metadata: TransactionMetadata): Promise<TransactionAction | undefined> {
    const transfers = this.moaSettingsService.getTransfers(metadata);
    if (!transfers) {
      return undefined;
    }

    const pair1Properties = transfers[0].properties;
    if (!pair1Properties) {
      return undefined;
    }

    const pair1Value = transfers[0].value;
    const pair1ValueDenominated = NumberUtils.toDenominatedString(pair1Value, pair1Properties.decimals);

    const numberOrArgumentsForOneSwap = 4;
    const numberOfSwaps = metadata.functionArgs.length / numberOrArgumentsForOneSwap;

    const swaps = [{
      properties: pair1Properties,
      value: pair1Value,
      denominatedValue: pair1ValueDenominated,
    }];

    for (let i = 0; i < numberOfSwaps; i++) {
      const tokenIdentifier = BinaryUtils.hexToString(metadata.functionArgs[i * numberOrArgumentsForOneSwap + 2]);
      const value = BinaryUtils.hexToBigInt(metadata.functionArgs[i * numberOrArgumentsForOneSwap + 3]);

      const pairProperties = await this.tokenTransferService.getTokenTransferProperties({ identifier: tokenIdentifier });
      if (!pairProperties) {
        return undefined;
      }

      const denominatedValue = NumberUtils.toDenominatedString(value, pairProperties.decimals);

      metadata.transfers?.push({
        value: value,
        properties: pairProperties,
      });

      swaps.push({
        denominatedValue,
        value,
        properties: pairProperties,
      });
    }

    const firstSwap = swaps[0];
    const lastSwap = swaps[swaps.length - 1];
    const intermediateSwaps = swaps.slice(1, swaps.length - 1);

    const description = `Swap ${firstSwap.denominatedValue} ${firstSwap.properties.ticker} for a minimum of ${lastSwap.denominatedValue} ${lastSwap.properties.ticker} with intermediate pair(s) ${intermediateSwaps.map(s => s.properties.ticker).join(', ')}`;

    return this.transactionActionDcdtNftRecognizerService.getMultiTransferAction(metadata, TransactionActionCategory.moa, 'multiSwap', description);
  }

  private getAddLiquidityAction(metadata: TransactionMetadata): TransactionAction | undefined {
    return this.transactionActionDcdtNftRecognizerService.getMultiTransferActionWithTicker(metadata, TransactionActionCategory.moa, MoaFunction.addLiquidity, 'Added liquidity for');
  }

  private getRemoveLiquidityAction(metadata: TransactionMetadata): TransactionAction | undefined {
    return this.transactionActionDcdtNftRecognizerService.getMultiTransferActionWithTicker(metadata, TransactionActionCategory.moa, MoaFunction.removeLiquidity, 'Removed liquidity with');
  }
}
