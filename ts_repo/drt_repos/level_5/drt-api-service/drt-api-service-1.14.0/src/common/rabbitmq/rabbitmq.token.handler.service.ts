import { Inject, Injectable } from '@nestjs/common';
import { CacheInfo } from '../../utils/cache.info';
import { NotifierEvent } from './entities/notifier.event';
import { DcdtService } from 'src/endpoints/dcdt/dcdt.service';
import { ClientProxy } from '@nestjs/microservices';
import { BinaryUtils, OriginLogger } from '@terradharitri/sdk-nestjs-common';
import { CacheService } from "@terradharitri/sdk-nestjs-cache";

@Injectable()
export class RabbitMqTokenHandlerService {
  private readonly logger = new OriginLogger(RabbitMqTokenHandlerService.name);

  constructor(
    private readonly cachingService: CacheService,
    private readonly dcdtService: DcdtService,
    @Inject('PUBSUB_SERVICE') private clientProxy: ClientProxy,
  ) { }

  public async handleTransferOwnershipEvent(event: NotifierEvent): Promise<boolean> {
    const tokenIdentifier = BinaryUtils.base64Decode(event.topics[0]);

    try {
      const dcdtProperties = await this.dcdtService.getDcdtTokenPropertiesRaw(tokenIdentifier);
      if (!dcdtProperties) {
        return false;
      }

      await this.invalidateKey(
        CacheInfo.DcdtProperties(tokenIdentifier).key,
        dcdtProperties,
        CacheInfo.DcdtProperties(tokenIdentifier).ttl
      );

      return true;
    } catch (error) {
      this.logger.error(`An unhandled error occurred when processing transferOwnership event for token with identifier '${tokenIdentifier}'`);
      this.logger.error(error);

      return false;
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
