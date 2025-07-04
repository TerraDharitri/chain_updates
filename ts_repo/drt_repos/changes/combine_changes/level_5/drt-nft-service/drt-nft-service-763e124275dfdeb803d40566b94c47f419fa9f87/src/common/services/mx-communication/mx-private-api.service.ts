import { HttpStatus, Injectable, Logger } from '@nestjs/common';
import { PerformanceProfiler } from 'src/modules/metrics/performance.profiler';
import { MetricsCollector } from 'src/modules/metrics/metrics.collector';
import * as Agent from 'agentkeepalive';
import { drtConfig } from 'src/config';
import { ApiNetworkProvider } from '@terradharitri/sdk-core';

@Injectable()
export class MxPrivateApiService {
  private privateApiProvider: ApiNetworkProvider;

  constructor(private readonly logger: Logger) {
    const keepAliveOptions = {
      maxSockets: drtConfig.keepAliveMaxSockets,
      maxFreeSockets: drtConfig.keepAliveMaxFreeSockets,
      timeout: parseInt(process.env.KEEPALIVE_TIMEOUT_DOWNSTREAM),
      freeSocketTimeout: drtConfig.keepAliveFreeSocketTimeout,
    };
    const httpAgent = new Agent(keepAliveOptions);
    const httpsAgent = new Agent.HttpsAgent(keepAliveOptions);

    this.privateApiProvider = new ApiNetworkProvider(process.env.NUMBAT_PRIVATE_API, {
      timeout: drtConfig.proxyTimeout,
      httpAgent: drtConfig.keepAlive ? httpAgent : null,
      httpsAgent: drtConfig.keepAlive ? httpsAgent : null,
      headers: {
        origin: 'NftService',
      },
      clientName: "nft-service",
    });
  }

  getPrivateService(): ApiNetworkProvider {
    return this.privateApiProvider;
  }

  async doPostGeneric(name: string, resourceUrl: string, payload: any): Promise<any> {
    try {
      const profiler = new PerformanceProfiler(`${name} ${resourceUrl}`);
      const service = this.getPrivateService();
      const response = await service.doPostGeneric(resourceUrl, payload);
      profiler.stop();

      MetricsCollector.setExternalCall(MxPrivateApiService.name, profiler.duration, name);

      return response;
    } catch (error) {
      if (error.inner?.response?.status === HttpStatus.NOT_FOUND) {
        return;
      }
      let customError = {
        method: 'POST',
        resourceUrl,
        response: error.inner?.response?.data,
        status: error.inner?.response?.status,
        message: error.message,
        name: error.name,
      };
      this.logger.error(`An error occurred while calling the drt private-api service on url ${resourceUrl}`, {
        path: `${MxPrivateApiService.name}.${this.doPostGeneric.name}`,
        error: customError,
      });
    }
  }
}
