/* eslint-disable @typescript-eslint/no-unused-vars */
import { Tokens, pairs } from 'src/modules/pair/mocks/pair.constants';
import { BaseDcdtToken, DcdtToken } from '../models/dcdtToken.model';
import { TokenService } from '../services/token.service';

export class TokenServiceMock {
    async tokenMetadata(tokenID: string): Promise<DcdtToken> {
        return Tokens(tokenID);
    }

    async baseTokenMetadata(tokenID: string): Promise<BaseDcdtToken> {
        return new BaseDcdtToken({
            identifier: Tokens(tokenID).identifier,
            decimals: Tokens(tokenID).decimals,
        });
    }

    async getDcdtTokenType(tokenID: string): Promise<string> {
        return Tokens(tokenID).type;
    }

    async getUniqueTokenIDs(activePool: boolean): Promise<string[]> {
        const tokenIDs = [];
        for (const pair of pairs) {
            tokenIDs.push(pair.firstToken.identifier);
            tokenIDs.push(pair.secondToken.identifier);
        }

        return [...new Set(tokenIDs)];
    }

    async getAllTokensMetadata(tokenIDs: string[]): Promise<DcdtToken[]> {
        return tokenIDs.map((tokenID) => Tokens(tokenID));
    }
    async getAllBaseTokensMetadata(
        tokenIDs: string[],
    ): Promise<BaseDcdtToken[]> {
        return tokenIDs.map(
            (tokenID) =>
                new BaseDcdtToken({
                    identifier: Tokens(tokenID).identifier,
                    decimals: Tokens(tokenID).decimals,
                }),
        );
    }
}

export const TokenServiceProvider = {
    provide: TokenService,
    useClass: TokenServiceMock,
};
