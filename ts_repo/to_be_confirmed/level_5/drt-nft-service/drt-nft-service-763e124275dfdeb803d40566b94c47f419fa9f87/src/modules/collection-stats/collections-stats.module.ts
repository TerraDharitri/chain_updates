import { Logger, Module } from '@nestjs/common';
import { DrtCommunicationModule } from 'src/common';
import { CollectionsStatsResolver } from './collections-stats.resolver';
import { CollectionsStatsService } from './collections-stats.service';

@Module({
  providers: [Logger, CollectionsStatsService, CollectionsStatsResolver],
  imports: [DrtCommunicationModule],
  exports: [CollectionsStatsService],
})
export class CollectionsStatsModuleGraph {}
