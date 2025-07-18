import { Injectable } from '@nestjs/common';
import { DrtApiService } from 'src/common';
import { Owner } from './models';

@Injectable()
export class OwnersService {
  constructor(private apiService: DrtApiService) {}

  async getOwnersForIdentifier(identifier: string, offset: number, limit: number): Promise<[Owner[], number]> {
    const [nfts, count] = await Promise.all([
      this.apiService.getOwnersForIdentifier(identifier, offset, limit),
      this.apiService.getOwnersForIdentifierCount(identifier),
    ]);
    const assets = nfts?.map((element) => Owner.fromApiOwner(element, identifier));
    return [assets, count];
  }
}
