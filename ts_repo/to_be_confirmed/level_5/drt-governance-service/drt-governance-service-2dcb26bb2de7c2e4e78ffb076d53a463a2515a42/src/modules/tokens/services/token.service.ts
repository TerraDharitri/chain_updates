import { Injectable } from '@nestjs/common';
import { DcdtToken } from '../models/dcdtToken.model';
import { CacheTtlInfo } from 'src/services/caching/cache.ttl.info';
import { DRTApiService } from 'src/services/dharitri-communication/drt.api.service';
import { CacheService } from '@terradharitri/sdk-nestjs-cache';

@Injectable()
export class TokenService {
    constructor(
        private readonly apiService: DRTApiService,
        protected readonly cachingService: CacheService,
    ) {}

    async getTokenMetadata(tokenID: string): Promise<DcdtToken> {
        if (tokenID === undefined) {
            return undefined;
        }
        const cacheKey = `token.${tokenID}`;
        const cachedToken = await this.cachingService.get<DcdtToken>(cacheKey);
        if (cachedToken && cachedToken !== undefined) {
            await this.cachingService.set<DcdtToken>(
                cacheKey,
                cachedToken,
                CacheTtlInfo.Token.remoteTtl,
                CacheTtlInfo.Token.localTtl,
            );
            return new DcdtToken(cachedToken);
        }

        const token = await this.apiService.getToken(tokenID);

        if (token !== undefined) {
            await this.cachingService.set<DcdtToken>(
                cacheKey,
                token,
                CacheTtlInfo.Token.remoteTtl,
                CacheTtlInfo.Token.localTtl,
            );

            return token;
        }

        return undefined;
    }
}
