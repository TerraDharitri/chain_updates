import { Injectable } from "@nestjs/common";
import { TransactionAction } from "../../entities/transaction.action";
import { TransactionMetadata } from "../../entities/transaction.metadata";
import { TransactionActionRecognizerInterface } from "../../transaction.action.recognizer.interface";
import { MoaFarmActionRecognizerService } from "./moa.farm.action.recognizer.service";
import { MoaPairActionRecognizerService } from "./moa.pair.action.recognizer.service";
import { MoaWrapActionRecognizerService } from "./moa.wrap.action.recognizer.service";
import { MoaDistributionActionRecognizerService } from "./moa.distribution.action.recognizer.service";
import { MoaLockedAssetActionRecognizerService } from "./moa.locked.asset.action.recognizer.service";
import { MoaSettingsService } from "../../../../moa/moa.settings.service";
import { ApiConfigService } from "src/common/api-config/api.config.service";

@Injectable()
export class TransactionActionMoaRecognizerService implements TransactionActionRecognizerInterface {
  constructor(
    private readonly pairActionRecognizer: MoaPairActionRecognizerService,
    private readonly farmActionRecognizer: MoaFarmActionRecognizerService,
    private readonly wrapActionRecognizer: MoaWrapActionRecognizerService,
    private readonly distributionRecognizer: MoaDistributionActionRecognizerService,
    private readonly lockedAssetRecognizer: MoaLockedAssetActionRecognizerService,
    private readonly moaSettingsService: MoaSettingsService,
    private readonly apiConfigService: ApiConfigService
  ) { }

  async isActive(): Promise<boolean> {
    const microServiceUrl = this.apiConfigService.getExchangeServiceUrl();
    if (!microServiceUrl) {
      return false;
    }

    const settings = await this.moaSettingsService.getSettings();
    return settings !== undefined;
  }

  async recognize(metadata: TransactionMetadata): Promise<TransactionAction | undefined> {
    const settings = await this.moaSettingsService.getSettings();
    if (!settings) {
      return undefined;
    }

    const isMoaInteraction = await this.moaSettingsService.isMoaInteraction(metadata);
    if (!isMoaInteraction) {
      return undefined;
    }

    return this.distributionRecognizer.recognize(settings, metadata) ??
      (await this.pairActionRecognizer.recognize(settings, metadata)) ??
      this.farmActionRecognizer.recognize(settings, metadata) ??
      this.wrapActionRecognizer.recognize(settings, metadata) ??
      this.lockedAssetRecognizer.recognize(settings, metadata);
  }
}
