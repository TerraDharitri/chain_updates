import { DcdtTokenPaymentType } from '../../attributes-decoder';

export type LockedFundsType = {
    funds: DcdtTokenPaymentType[];
    lockedEpoch: number;
};

export type EscrowBaseEventType = {
    sender: string;
    receiver: string;
    lockedFunds: LockedFundsType;
};
