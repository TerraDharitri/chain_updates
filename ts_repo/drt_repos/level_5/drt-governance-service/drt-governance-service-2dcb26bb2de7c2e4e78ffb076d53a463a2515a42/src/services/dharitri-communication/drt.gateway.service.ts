import { Inject, Injectable } from '@nestjs/common';
import { WINSTON_MODULE_PROVIDER } from 'nest-winston';
import { drtConfig } from 'src/config';
import { ApiConfigService } from 'src/helpers/api.config.service';
import { Logger } from 'winston';
import Agent, { HttpsAgent } from 'agentkeepalive';
import axios, { AxiosRequestConfig } from 'axios';
import { PerformanceProfiler } from 'src/utils/performance.profiler';
import { MetricsCollector } from 'src/utils/metrics.collector';
import { Address } from '@terradharitri/sdk-core/out';

@Injectable()
export class DRTGatewayService {
    private url: string;
    private config: AxiosRequestConfig;

    constructor(
        private readonly apiConfigService: ApiConfigService,
        @Inject(WINSTON_MODULE_PROVIDER) private readonly logger: Logger,
    ) {
        const keepAliveOptions = {
            maxSockets: drtConfig.keepAliveMaxSockets,
            maxFreeSockets: drtConfig.keepAliveMaxFreeSockets,
            timeout: this.apiConfigService.getKeepAliveTimeoutDownstream(),
            freeSocketTimeout: drtConfig.keepAliveFreeSocketTimeout,
            keepAlive: true,
        };
        const httpAgent = new Agent(keepAliveOptions);
        const httpsAgent = new HttpsAgent(keepAliveOptions);
        this.url = process.env.DRT_GATEWAY_URL;

        this.config = {
            timeout: drtConfig.proxyTimeout,
            httpAgent: drtConfig.keepAlive ? httpAgent : null,
            httpsAgent: drtConfig.keepAlive ? httpsAgent : null,
        };
    }

    async getSCStorageKey(address: string, key: string): Promise<any> {
        return await this.doGetGeneric(
            this.getSCStorageKey.name,
            `address/${address}/key/${Buffer.from(key).toString('hex')}`,
            (response) => response.data.value,
        );
    }

    async getSCStorageKeys(address: string, keys: any[]): Promise<any> {
        let fullKey = '';
        for (const key of keys) {
            switch (typeof key) {
                case 'number':
                    fullKey = fullKey.concat(key.toString(16).padStart(8, '0'));
                    break;
                case 'string':
                    fullKey = fullKey.concat(Buffer.from(key).toString('hex'));
                    break;
                case 'object':
                    if (key instanceof Address) {
                        fullKey = fullKey.concat(key.hex());
                        break;
                    }
                    break;
            }
        }

        const uri =
            fullKey !== ''
                ? `address/${address}/key/${fullKey}`
                : `address/${address}/keys`;

        const response = await this.doGetGeneric(
            this.getSCStorageKey.name,
            uri,
            (response) => response.data,
        );
        return fullKey !== '' ? response.value : response.pairs;
    }

    /**
     * Get method that receives the resource url and on callback the method used to map the response.
     */
    private async doGetGeneric(
        name: string,
        resourceUrl: string,
        callback: (response: any) => any,
    ): Promise<any> {
        const response = await this.doGet(name, resourceUrl);
        return callback(response);
    }

    private async doGet(name: string, resourceUrl: string): Promise<any> {
        const profiler = new PerformanceProfiler(`${name} ${resourceUrl}`);
        try {
            const url = `${this.url}/${resourceUrl}`;
            const response = await axios.get(url, this.config);
            return response.data;
        } catch (error) {
            this.handleApiError(error, resourceUrl);
        } finally {
            profiler.stop();

            MetricsCollector.setExternalCall(
                DRTGatewayService.name,
                name,
                profiler.duration,
            );
        }
    }

    private handleApiError(error: any, resourceUrl: string) {
        if (!error.response) {
            this.logger.warn(error);
            throw new Error(`Cannot GET ${resourceUrl}: [${error}]`);
        }

        const errorData = error.response.data;
        const originalErrorMessage =
            errorData.error || errorData.message || JSON.stringify(errorData);
        throw new Error(`Cannot GET ${resourceUrl}: [${originalErrorMessage}]`);
    }
}
