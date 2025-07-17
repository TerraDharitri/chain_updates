import { DcdtTokenPaymentType } from '../../../attributes-decoder';
import { FarmTokenAttributesTypeV2 } from '../../../attributes-decoder/farm/farm.token.types';

export type FarmEventsTopicsTypeV2 = {
    eventName: string;
    caller: string;
    epoch: number;
    block: number;
    timestamp: number;
    farmingTokenID: string;
};

export type EnterFarmEventTypeV2 = {
    farmingToken: DcdtTokenPaymentType;
    farmToken: DcdtTokenPaymentType;
    farmSupply: string;
    rewardTokenID: string;
    rewardTokenReserves: string;
    createdWithMerge: boolean;
    farmAttributes: FarmTokenAttributesTypeV2;
};

export type ExitFarmEventTypeV2 = {
    farmingToken: DcdtTokenPaymentType;
    farmToken: DcdtTokenPaymentType;
    farmSupply: string;
    rewardTokens: DcdtTokenPaymentType;
    rewardTokenReserves: string;
    farmAttributes: FarmTokenAttributesTypeV2;
};

export type ClaimRewardsEventTypeV2 = {
    oldFarmToken: DcdtTokenPaymentType;
    newFarmToken: DcdtTokenPaymentType;
    farmSupply: string;
    rewardTokens: DcdtTokenPaymentType;
    rewardTokenReserves: string;
    oldFarmAttributes: FarmTokenAttributesTypeV2;
    newFarmAttributes: FarmTokenAttributesTypeV2;
    createdWithMerge: boolean;
};
