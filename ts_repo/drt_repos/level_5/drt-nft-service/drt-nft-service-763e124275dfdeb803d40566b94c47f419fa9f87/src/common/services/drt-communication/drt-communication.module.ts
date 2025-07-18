import { Logger, Module } from '@nestjs/common';
import { ApiService } from './api.service';
import { DrtApiService } from './drt-api.service';
import { DrtElasticService } from './drt-elastic.service';
import { DrtFeedService } from './drt-feed.service';
import { DrtIdentityService } from './drt-identity.service';
import { DrtPrivateApiService } from './drt-private-api.service';
import { DrtProxyService } from './drt-proxy.service';
import { DrtStatsService } from './drt-stats.service';
import { SlackReportService } from './slack-report.service';
import { ConfigService } from '@nestjs/config';
import { ApiConfigService } from 'src/modules/common/api-config/api.config.service';
import { DrtDataApiService } from './drt-data.service';
import { DrtToolsService } from './drt-tools.service';

@Module({
  providers: [
    Logger,
    ApiService,
    ConfigService,
    ApiConfigService,
    DrtProxyService,
    DrtApiService,
    DrtPrivateApiService,
    DrtStatsService,
    DrtElasticService,
    DrtIdentityService,
    DrtFeedService,
    SlackReportService,
    DrtDataApiService,
    DrtToolsService,
  ],
  exports: [
    ApiService,
    DrtProxyService,
    DrtStatsService,
    DrtElasticService,
    DrtApiService,
    DrtPrivateApiService,
    DrtIdentityService,
    DrtFeedService,
    SlackReportService,
    DrtDataApiService,
    DrtToolsService,
  ],
})
export class DrtCommunicationModule {}
