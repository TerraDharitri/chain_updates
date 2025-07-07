import { forwardRef, Module } from '@nestjs/common';
import { MongooseModule } from '@nestjs/mongoose';
import { PairModule } from '../pair/pair.module';
import { RouterModule } from '../router/router.module';
import { DcdtTokenDbModel, DcdtTokenSchema } from './schemas/token.schema';
import { TokenRepositoryService } from './services/token.repository.service';
import { TokenService } from './services/token.service';
import { TokensResolver } from './token.resolver';
import { DatabaseModule } from 'src/services/database/database.module';
import { TokenComputeService } from './services/token.compute.service';
import { TokenSetterService } from './services/token.setter.service';
import { MXCommunicationModule } from 'src/services/dharitri-communication/drt.communication.module';
import { NftTokenResolver } from './nftToken.resolver';
import { AnalyticsModule } from 'src/services/analytics/analytics.module';
import { TokenFilteringService } from './services/token.filtering.service';
import { ElasticSearchModule } from 'src/services/elastic-search/elastic.search.module';
import { TokenLoader } from './services/token.loader';

@Module({
    imports: [
        MXCommunicationModule,
        forwardRef(() => PairModule),
        forwardRef(() => RouterModule),
        DatabaseModule,
        MongooseModule.forFeature([
            { name: DcdtTokenDbModel.name, schema: DcdtTokenSchema },
        ]),
        AnalyticsModule,
        ElasticSearchModule,
    ],
    providers: [
        TokenLoader,
        TokenService,
        TokenSetterService,
        TokenComputeService,
        TokenRepositoryService,
        TokensResolver,
        NftTokenResolver,
        TokenFilteringService,
    ],
    exports: [
        TokenRepositoryService,
        TokenLoader,
        TokenService,
        TokenSetterService,
        TokenComputeService,
        TokenFilteringService,
    ],
})
export class TokenModule {}
