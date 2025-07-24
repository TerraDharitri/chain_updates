import { forwardRef, Module } from '@nestjs/common';
import { CampaignsQueriesResolver } from './campaigns-queries.resolver';
import { DrtCommunicationModule } from 'src/common';
import { CampaignsMutationsResolver } from './campaigns-mutations.resolver';
import { NftMinterAbiService } from './nft-minter.abi.service';
import { CampaignsService } from './campaigns.service';
import { PubSubListenerModule } from 'src/pubsub/pub.sub.listener.module';
import { CommonModule } from 'src/common.module';
import { AuthModule } from '../auth/auth.module';
import { CampaignsCachingService } from './campaigns-caching.service';

@Module({
  providers: [CampaignsMutationsResolver, CampaignsQueriesResolver, NftMinterAbiService, CampaignsService, CampaignsCachingService],
  imports: [PubSubListenerModule, DrtCommunicationModule, CommonModule, forwardRef(() => AuthModule)],
  exports: [CampaignsService],
})
export class CampaignsModuleGraph {}
