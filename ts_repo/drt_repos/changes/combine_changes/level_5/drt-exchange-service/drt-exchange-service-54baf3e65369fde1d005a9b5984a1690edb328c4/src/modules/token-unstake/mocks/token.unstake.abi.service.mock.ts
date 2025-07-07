import { DcdtTokenPaymentModel } from 'src/modules/tokens/models/dcdt.token.payment.model';
import { UnstakePairModel } from '../models/token.unstake.model';
import { ITokenUnstakeAbiService } from '../services/interfaces';
import { TokenUnstakeAbiService } from '../services/token.unstake.abi.service';

export class TokenUnstakeAbiServiceMock implements ITokenUnstakeAbiService {
    unbondEpochs(): Promise<number> {
        throw new Error('Method not implemented.');
    }
    feesBurnPercentage(): Promise<number> {
        throw new Error('Method not implemented.');
    }
    feesCollectorAddress(): Promise<string> {
        throw new Error('Method not implemented.');
    }
    energyFactoryAddress(): Promise<string> {
        throw new Error('Method not implemented.');
    }

    async unlockedTokensForUser(
        userAddress: string,
    ): Promise<UnstakePairModel[]> {
        return [
            new UnstakePairModel({
                lockedTokens: new DcdtTokenPaymentModel({
                    amount: '1000000000000000000',
                    tokenIdentifier: 'XMOA-123456',
                    tokenNonce: 2,
                }),
                unlockedTokens: new DcdtTokenPaymentModel({
                    amount: '1000000000000000000',
                    tokenIdentifier: 'XMOA-123456',
                    tokenNonce: 1,
                }),
                unlockEpoch: 1,
            }),
        ];
    }
}

export const TokenUnstakeAbiServiceProvider = {
    provide: TokenUnstakeAbiService,
    useClass: TokenUnstakeAbiServiceMock,
};
