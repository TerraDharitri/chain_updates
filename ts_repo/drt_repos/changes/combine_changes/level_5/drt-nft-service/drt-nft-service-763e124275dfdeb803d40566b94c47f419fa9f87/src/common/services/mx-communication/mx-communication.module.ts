import { Logger, Module } from '@nestjs/common';
import { ApiService } from './api.service';
import { MxApiService } from './drt-api.service';
import { MxElasticService } from './drt-elastic.service';
import { MxFeedService } from './drt-feed.service';
import { MxIdentityService } from './drt-identity.service';
import { MxPrivateApiService } from './drt-private-api.service';
import { MxProxyService } from './drt-proxy.service';
import { MxStatsService } from './drt-stats.service';
import { SlackReportService } from './slack-report.service';
import { ConfigService } from '@nestjs/config';
import { ApiConfigService } from 'src/modules/common/api-config/api.config.service';
import { MxDataApiService } from './drt-data.service';
import { MxToolsService } from './drt-tools.service';

@Module({
  providers: [
    Logger,
    ApiService,
    ConfigService,
    ApiConfigService,
    MxProxyService,
    MxApiService,
    MxPrivateApiService,
    MxStatsService,
    MxElasticService,
    MxIdentityService,
    MxFeedService,
    SlackReportService,
    MxDataApiService,
    MxToolsService,
  ],
  exports: [
    ApiService,
    MxProxyService,
    MxStatsService,
    MxElasticService,
    MxApiService,
    MxPrivateApiService,
    MxIdentityService,
    MxFeedService,
    SlackReportService,
    MxDataApiService,
    MxToolsService,
  ],
})
export class MxCommunicationModule {}
