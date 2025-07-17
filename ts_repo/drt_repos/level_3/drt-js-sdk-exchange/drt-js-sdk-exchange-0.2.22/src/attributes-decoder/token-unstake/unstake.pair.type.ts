import { DcdtTokenPaymentType } from '../dcdt-token-payment';

export type UnstakePairType = {
    unlockEpoch: number;
    lockedTokens: DcdtTokenPaymentType;
    unlockedTokens: DcdtTokenPaymentType;
};
