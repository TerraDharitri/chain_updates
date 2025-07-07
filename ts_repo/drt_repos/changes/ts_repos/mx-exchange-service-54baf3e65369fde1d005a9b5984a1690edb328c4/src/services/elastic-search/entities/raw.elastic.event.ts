import { RawEventType } from '@terradharitri/sdk-exchange';

export type RawElasticEventType = RawEventType & {
    logAddress: string;
    shardID: number;
    timestamp: number;
    txOrder: number;
    txHash: string;
    order: number;
};
