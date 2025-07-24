import { Inject, Injectable } from '@nestjs/common';
import { WINSTON_MODULE_PROVIDER } from 'nest-winston';
import { Constants } from '@terradharitri/sdk-nestjs-common';
import { Logger } from 'winston';
import { CacheService } from '@terradharitri/sdk-nestjs-cache';
import { DRTApiService } from '../dharitri-communication/drt.api.service';
import { GenericGetterService } from '../generics/generic.getter.service';

@Injectable()
export class ContextGetterService extends GenericGetterService {
    constructor(
        protected readonly cachingService: CacheService,
        @Inject(WINSTON_MODULE_PROVIDER) protected readonly logger: Logger,
        private readonly apiService: DRTApiService,
    ) {
        super(cachingService, logger);
        this.baseKey = 'context';
    }

    async getCurrentEpoch(): Promise<number> {
        const cacheKey = this.getCacheKey('currentEpoch');
        return await this.getData(
            cacheKey,
            async () => (await this.apiService.getStats()).epoch,
            Constants.oneMinute(),
        );
    }
}
