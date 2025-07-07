import { DcdtTokenPaymentType } from '../dcdt-token-payment/dcdt.token.payment.type';

export type WrappedLpAttributesType = {
    lpTokenID: string | undefined;
    lpTokenTotalAmount: string | undefined;
    lockedAssetsInvested: string | undefined;
    lockedAssetsNonce: number | undefined;
};

export type WrappedFarmAttributesType = {
    farmTokenID: string | undefined;
    farmTokenNonce: number | undefined;
    farmTokenAmount: string | undefined;
    farmingTokenID: string | undefined;
    farmingTokenNonce: number | undefined;
    farmingTokenAmount: string | undefined;
};

export type WrappedLpTokenAttributesTypeV2 = {
    lpTokenID: string;
    lpTokenAmount: string;
    lockedTokens: DcdtTokenPaymentType;
};

export type WrappedFarmTokenAttributesTypeV2 = {
    farmToken: DcdtTokenPaymentType;
    proxyFarmingToken: DcdtTokenPaymentType;
};
