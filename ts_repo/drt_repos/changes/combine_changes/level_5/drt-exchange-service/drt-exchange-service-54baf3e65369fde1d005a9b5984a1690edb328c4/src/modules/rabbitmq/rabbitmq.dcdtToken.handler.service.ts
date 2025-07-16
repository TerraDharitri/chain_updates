import {
    DcdtLocalBurnEvent,
    DcdtLocalMintEvent,
} from '@terradharitri/sdk-exchange';
import { Inject, Injectable } from '@nestjs/common';
import { RedisPubSub } from 'graphql-redis-subscriptions';
import { WINSTON_MODULE_PROVIDER } from 'nest-winston';
import { Constants } from '@terradharitri/sdk-nestjs-common';
import { CacheService } from 'src/services/caching/cache.service';
import { MXApiService } from 'src/services/dharitri-communication/drt.api.service';
import { PUB_SUB } from 'src/services/redis.pubSub.module';
import { generateCacheKeyFromParams } from 'src/utils/generate-cache-key';
import { Logger } from 'winston';

@Injectable()
export class RabbitMQDcdtTokenHandlerService {
    private invalidatedKeys = [];

    constructor(
        private readonly apiService: MXApiService,
        private readonly cachingService: CacheService,
        @Inject(PUB_SUB) private pubSub: RedisPubSub,
        @Inject(WINSTON_MODULE_PROVIDER) private readonly logger: Logger,
    ) {}

    async handleDcdtTokenEvent(
        event: DcdtLocalMintEvent | DcdtLocalBurnEvent,
    ): Promise<void> {
        const cacheKey = generateCacheKeyFromParams(
            'context',
            event.getTopics().tokenID,
        );
        const token = await this.apiService.getToken(event.getTopics().tokenID);
        await this.cachingService.set(cacheKey, token, Constants.oneHour());
        this.invalidatedKeys.push(cacheKey);
        await this.deleteCacheKeys();
    }

    private async deleteCacheKeys() {
        await this.pubSub.publish('deleteCacheKeys', this.invalidatedKeys);
        this.invalidatedKeys = [];
    }
}
