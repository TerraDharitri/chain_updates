import { HitResponse } from './models/elastic-search';

export interface AddressTransactionCount {
  contractAddress: string;
  transactionCount: number;
}

export class DrtElasticServiceMock {
  async getNftHistory(
    collection: string,
    nonce: string,
    size: number,
    timestamp: number | string,
  ): Promise<[HitResponse[], number, number]> {
    return [[], 0, 0];
  }
}
