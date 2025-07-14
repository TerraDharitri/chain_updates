import { Injectable } from '@nestjs/common';
import { PaginationArgs } from 'src/modules/dex.model';
import { PairService } from 'src/modules/pair/services/pair.service';
import { IDcdtToken } from 'src/modules/tokens/models/dcdtToken.interface';
import {
    DcdtToken,
    DcdtTokenType,
} from 'src/modules/tokens/models/dcdtToken.model';
import { TokenService } from 'src/modules/tokens/services/token.service';
import { MXApiService } from 'src/services/dharitri-communication/drt.api.service';
import { UserToken } from '../models/user.model';
import { UserDcdtComputeService } from './dcdt.compute.service';
import { RouterAbiService } from 'src/modules/router/services/router.abi.service';

@Injectable()
export class UserDcdtService {
    constructor(
        private readonly apiService: MXApiService,
        private readonly pairService: PairService,
        private readonly tokenService: TokenService,
        private readonly routerAbi: RouterAbiService,
        private readonly userDcdtCompute: UserDcdtComputeService,
    ) {}

    async getAllDcdtTokens(
        userAddress: string,
        pagination: PaginationArgs,
        inputTokens?: IDcdtToken[],
    ): Promise<UserToken[]> {
        const [userTokens, uniquePairTokens, pairsAddresses] =
            await Promise.all([
                inputTokens
                    ? Promise.resolve(inputTokens)
                    : this.apiService.getTokensForUser(
                          userAddress,
                          pagination.offset,
                          pagination.limit,
                      ),
                this.tokenService.getUniqueTokenIDs(false),
                this.routerAbi.pairsAddress(),
            ]);
        const lpTokensIDs = await this.pairService.getAllLpTokensIds(
            pairsAddresses,
        );

        const userPairDcdtTokens = userTokens
            .filter((token) => uniquePairTokens.includes(token.identifier))
            .map((token) => new DcdtToken(token));

        const userLpDcdtTokens = userTokens
            .filter((token) => lpTokensIDs.includes(token.identifier))
            .map((token) => new DcdtToken(token));

        const orderedPairAddresses = [];
        userLpDcdtTokens.forEach((token, index) => {
            userLpDcdtTokens[index].type = DcdtTokenType.FungibleLpToken;

            const lpTokenIndex = lpTokensIDs.findIndex(
                (elem) => elem === token.identifier,
            );
            orderedPairAddresses.push(pairsAddresses[lpTokenIndex]);
        });

        const detailedLpTokens = await this.userDcdtCompute.allLpTokensUSD(
            userLpDcdtTokens,
            orderedPairAddresses,
        );
        const detailedDcdtTokens = await this.userDcdtCompute.allDcdtTokensUSD(
            userPairDcdtTokens,
        );

        return [...detailedDcdtTokens, ...detailedLpTokens];
    }
}
