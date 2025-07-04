import { DcdtToken } from 'src/modules/tokens/models/dcdtToken.model';
import { NftToken } from 'src/modules/tokens/models/nftToken.model';
import { Tokens } from 'src/modules/pair/mocks/pair.constants';
import { MXApiService } from './drt.api.service';
import { Address } from '@terradharitri/sdk-core/out';

export class MXApiServiceMock {
    async getCurrentEpoch(): Promise<number> {
        return 1;
    }

    async getToken(tokenID: string): Promise<DcdtToken> {
        return Tokens(tokenID);
    }

    async getAccountStats(address: string): Promise<any | undefined> {
        return {
            address: 'user_address_1',
            nonce: 1,
            balance: '1000000000000000000',
            rootHash: 'Xh9naslK9HMLbd3O6ueH04pdfDO0TIsXoM9BF9jouzs=',
            txCount: 1,
            shard: 1,
        };
    }

    async getTokensForUser(address: string): Promise<DcdtToken[]> {
        return [
            new DcdtToken({
                ...Tokens('MOA-123456'),
                balance: '1000000000000000000',
            }),
        ];
    }

    async getNftsCountForUser(address: string): Promise<number> {
        return 1;
    }

    async getNftsForUser(address: string): Promise<NftToken[]> {
        return [
            {
                collection: 'REWAMOAFL-abcdef',
                ticker: 'REWAMOAFL',
                name: 'FarmToken',
                type: 'SemiFungibleDCDT',
                decimals: 18,
                balance: '1000000000000000000',
                identifier: 'REWAMOAFL-abcdef-01',
                attributes: 'AAAABQeMCWDbAAAAAAAAAF8CAQ==',
                creator: Address.fromHex(
                    '0000000000000000000000000000000000000000000000000000000000000021',
                ).bech32(),
                nonce: 1,
                royalties: 0,
                timestamp: 0,
                uris: [],
                url: '',
                tags: [],
            },
        ];
    }

    getNftAttributesByTokenIdentifier(
        address: string,
        nftIdentifier: string,
    ): Promise<string> {
        throw new Error('Method not implemented.');
    }
}

export const MXApiServiceProvider = {
    provide: MXApiService,
    useClass: MXApiServiceMock,
};
