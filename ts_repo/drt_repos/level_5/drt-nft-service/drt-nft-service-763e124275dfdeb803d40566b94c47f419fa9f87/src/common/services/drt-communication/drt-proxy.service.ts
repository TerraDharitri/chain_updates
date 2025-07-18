import { ProxyNetworkProvider } from '@terradharitri/sdk-core';
import { Injectable } from '@nestjs/common';
import * as Agent from 'agentkeepalive';
import { drtConfig } from '../../../config';

@Injectable()
export class MxProxyService {
  private readonly proxy: ProxyNetworkProvider;
  constructor() {
    const keepAliveOptions = {
      maxSockets: drtConfig.keepAliveMaxSockets,
      maxFreeSockets: drtConfig.keepAliveMaxFreeSockets,
      timeout: parseInt(process.env.KEEPALIVE_TIMEOUT_DOWNSTREAM),
      freeSocketTimeout: drtConfig.keepAliveFreeSocketTimeout,
    };
    const httpAgent = new Agent(keepAliveOptions);
    const httpsAgent = new Agent.HttpsAgent(keepAliveOptions);

    this.proxy = new ProxyNetworkProvider(process.env.NUMBAT_GATEWAY, {
      timeout: parseInt(process.env.KEEPALIVE_TIMEOUT_DOWNSTREAM),
      httpAgent: drtConfig.keepAlive ? httpAgent : null,
      httpsAgent: drtConfig.keepAlive ? httpsAgent : null,
      clientName: 'nft-service',
    });
  }

  getService(): ProxyNetworkProvider {
    return this.proxy;
  }
}
