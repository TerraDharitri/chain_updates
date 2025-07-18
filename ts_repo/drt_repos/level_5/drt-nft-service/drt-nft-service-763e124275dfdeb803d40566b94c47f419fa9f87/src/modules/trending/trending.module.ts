import { Logger, Module } from '@nestjs/common';
import { DrtCommunicationModule } from 'src/common';
import { TrendingResolver } from './trending.resolver';
import { TrendingService } from './trending.service';

@Module({
  providers: [Logger, TrendingService, TrendingResolver],
  imports: [DrtCommunicationModule],
})
export class TrendingModuleGraph {}
