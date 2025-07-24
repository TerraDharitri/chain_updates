import { Inject, Injectable } from "@nestjs/common";
import { Cron, CronExpression } from "@nestjs/schedule";
import { ClientProxy } from "@nestjs/microservices";
import { CacheInfo } from "src/utils/cache.info";
import { MoaSettingsService } from "src/endpoints/moa/moa.settings.service";
import { MoaEconomicsService } from "src/endpoints/moa/moa.economics.service";
import { MoaPairService } from "src/endpoints/moa/moa.pair.service";
import { MoaTokenService } from "src/endpoints/moa/moa.token.service";
import { MoaFarmService } from "src/endpoints/moa/moa.farm.service";
import { Lock, Locker } from "@terradharitri/sdk-nestjs-common";
import { CacheService } from "@terradharitri/sdk-nestjs-cache";

@Injectable()
export class MoaWarmerService {
  constructor(
    private readonly cachingService: CacheService,
    @Inject('PUBSUB_SERVICE') private clientProxy: ClientProxy,
    private readonly moaEconomicsService: MoaEconomicsService,
    private readonly drtPairsService: MoaPairService,
    private readonly moaTokensService: MoaTokenService,
    private readonly moaSettingsService: MoaSettingsService,
    private readonly moaFarmsService: MoaFarmService,
  ) { }

  @Cron(CronExpression.EVERY_MINUTE)
  async handleMoaInvalidations() {
    await Locker.lock('Refreshing moa pairs', async () => {
      await this.drtPairsService.refreshMoaPairs();
    }, true);

    await Locker.lock('Refreshing moa economics', async () => {
      await this.moaEconomicsService.refreshMoaEconomics();
    }, true);

    await Locker.lock('Refreshing moa tokens', async () => {
      await this.moaTokensService.refreshMoaTokens();
    }, true);

    await Locker.lock('Refreshing moa farms', async () => {
      await this.moaFarmsService.refreshMoaFarms();
    }, true);

    await Locker.lock('Refreshing moa settings', async () => {
      await this.moaSettingsService.refreshSettings();
    }, true);
  }

  @Cron(CronExpression.EVERY_10_MINUTES)
  @Lock({ name: 'Moa settings invalidations' })
  async handleMoaSettings() {
    const settings = await this.moaSettingsService.getSettingsRaw();
    if (settings) {
      await this.invalidateKey(CacheInfo.MoaSettings.key, settings, CacheInfo.MoaSettings.ttl);
    }
  }

  private async invalidateKey(key: string, data: any, ttl: number) {
    await this.cachingService.set(key, data, ttl);
    this.refreshCacheKey(key, ttl);
  }

  private refreshCacheKey(key: string, ttl: number) {
    this.clientProxy.emit('refreshCacheKey', { key, ttl });
  }
}
