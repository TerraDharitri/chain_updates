import { Injectable } from '@nestjs/common';
import { PairService } from 'src/modules/pair/services/pair.service';
import { DcdtToken } from 'src/modules/tokens/models/dcdtToken.model';
import { computeValueUSD } from 'src/utils/token.converters';
import { UserToken } from '../models/user.model';
import { PairComputeService } from 'src/modules/pair/services/pair.compute.service';
import { TokenComputeService } from 'src/modules/tokens/services/token.compute.service';

@Injectable()
export class UserDcdtComputeService {
    constructor(
        private readonly pairService: PairService,
        private readonly pairCompute: PairComputeService,
        private readonly tokenCompute: TokenComputeService,
    ) {}

    async dcdtTokenUSD(dcdtToken: DcdtToken): Promise<UserToken> {
        const tokenPriceUSD = await this.pairCompute.tokenPriceUSD(
            dcdtToken.identifier,
        );
        return new UserToken({
            ...dcdtToken,
            valueUSD: computeValueUSD(
                dcdtToken.balance,
                dcdtToken.decimals,
                tokenPriceUSD,
            ).toFixed(),
        });
    }

    async allDcdtTokensUSD(dcdtTokens: DcdtToken[]): Promise<UserToken[]> {
        const allTokenPrices =
            await this.tokenCompute.getAllTokensPriceDerivedUSD(
                dcdtTokens.map((dcdtToken) => dcdtToken.identifier),
            );

        return dcdtTokens.map(
            (dcdtToken, index) =>
                new UserToken({
                    ...dcdtToken,
                    valueUSD: computeValueUSD(
                        dcdtToken.balance,
                        dcdtToken.decimals,
                        allTokenPrices[index],
                    ).toFixed(),
                }),
        );
    }

    async lpTokenUSD(
        dcdtToken: DcdtToken,
        pairAddress: string,
    ): Promise<UserToken> {
        const valueUSD = await this.pairService.getLiquidityPositionUSD(
            pairAddress,
            dcdtToken.balance,
        );
        return new UserToken({
            ...dcdtToken,
            valueUSD: valueUSD,
            pairAddress,
        });
    }

    async allLpTokensUSD(
        dcdtTokens: DcdtToken[],
        pairAddresses: string[],
    ): Promise<UserToken[]> {
        const valuesUSD = await this.pairService.getAllLiquidityPositionsUSD(
            pairAddresses,
            dcdtTokens.map((token) => token.balance),
        );

        return dcdtTokens.map(
            (dcdtToken, index) =>
                new UserToken({
                    ...dcdtToken,
                    valueUSD: valuesUSD[index],
                    pairAddress: pairAddresses[index],
                }),
        );
    }
}
