import { Test, TestingModule } from '@nestjs/testing';
import { TokenComputeService } from '../services/token.compute.service';
import { PairAbiServiceProvider } from 'src/modules/pair/mocks/pair.abi.service.mock';
import { PairComputeServiceProvider } from 'src/modules/pair/mocks/pair.compute.service.mock';
import { RouterAbiServiceProvider } from 'src/modules/router/mocks/router.abi.service.mock';
import { MXDataApiServiceProvider } from 'src/services/dharitri-communication/drt.data.api.service.mock';
import { PairService } from 'src/modules/pair/services/pair.service';
import { WrapAbiServiceProvider } from 'src/modules/wrapping/mocks/wrap.abi.service.mock';
import { TokenServiceProvider } from '../mocks/token.service.mock';
import { ContextGetterServiceProvider } from 'src/services/context/mocks/context.getter.service.mock';
import { tokenProviderUSD } from 'src/config';
import { ConfigModule } from '@nestjs/config';
import { WinstonModule } from 'nest-winston';
import { ApiConfigService } from 'src/helpers/api.config.service';
import winston from 'winston';
import { DynamicModuleUtils } from 'src/utils/dynamic.module.utils';
import { AnalyticsQueryServiceProvider } from 'src/services/analytics/mocks/analytics.query.service.mock';
import { MXApiServiceProvider } from 'src/services/dharitri-communication/drt.api.service.mock';
import { ElasticSearchModule } from 'src/services/elastic-search/elastic.search.module';

describe('TokenComputeService', () => {
    let module: TestingModule;

    beforeAll(async () => {
        module = await Test.createTestingModule({
            imports: [
                WinstonModule.forRoot({
                    transports: [new winston.transports.Console({})],
                }),
                ConfigModule.forRoot({}),
                DynamicModuleUtils.getCacheModule(),
                ElasticSearchModule,
            ],
            providers: [
                PairAbiServiceProvider,
                PairComputeServiceProvider,
                PairService,
                WrapAbiServiceProvider,
                TokenServiceProvider,
                RouterAbiServiceProvider,
                MXDataApiServiceProvider,
                ContextGetterServiceProvider,
                TokenComputeService,
                ApiConfigService,
                AnalyticsQueryServiceProvider,
                MXApiServiceProvider,
            ],
        }).compile();
    });

    it('should be defined', () => {
        const service: TokenComputeService =
            module.get<TokenComputeService>(TokenComputeService);
        expect(service).toBeDefined();
    });

    it('should compute token price derived REWA for tokenProviderUSD', async () => {
        const service: TokenComputeService =
            module.get<TokenComputeService>(TokenComputeService);
        const price = await service.computeTokenPriceDerivedREWA(
            tokenProviderUSD,
            [],
        );
        expect(price).toEqual('1');
    });

    it('should compute token price derived REWA for MOA-123456', async () => {
        const service: TokenComputeService =
            module.get<TokenComputeService>(TokenComputeService);
        const price = await service.computeTokenPriceDerivedREWA(
            'MOA-123456',
            [],
        );
        expect(price).toEqual('0.001');
    });

    it('should compute token price derived REWA for TOK4-123456', async () => {
        const service: TokenComputeService =
            module.get<TokenComputeService>(TokenComputeService);
        const price = await service.computeTokenPriceDerivedREWA(
            'TOK4-123456',
            [],
        );
        expect(price).toEqual('0.01');
    });

    it('should compute token price derived REWA for TOK5-123456', async () => {
        const service: TokenComputeService =
            module.get<TokenComputeService>(TokenComputeService);
        const price = await service.computeTokenPriceDerivedREWA(
            'TOK5-123456',
            [],
        );
        expect(price).toEqual('0.01010101010101010101010101');
    });

    it('should compute token price derived REWA for TOK6-123456', async () => {
        const service: TokenComputeService =
            module.get<TokenComputeService>(TokenComputeService);
        const price = await service.computeTokenPriceDerivedREWA(
            'TOK6-123456',
            [],
        );
        expect(price).toEqual('0.01010101010101010101010101');
    });
});
