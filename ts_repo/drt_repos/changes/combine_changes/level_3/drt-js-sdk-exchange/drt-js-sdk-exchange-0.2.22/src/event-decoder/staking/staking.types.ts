import { DcdtTokenPaymentType } from '../../attributes-decoder';
import { StakingFarmTokenAttributesType } from '../../attributes-decoder/staking/staking.farm.token.types';

export type StakingEventsTopicsType = {
    eventName: string;
    caller: string;
    epoch: number;
    block: number;
    timestamp: number;
    farmingTokenID: string;
};

export type StakeEventType = {
    farmingToken: DcdtTokenPaymentType;
    farmToken: DcdtTokenPaymentType;
    farmSupply: string;
    rewardTokenID: string;
    rewardTokenReserves: string;
    createdWithMerge: boolean;
    farmAttributes: StakingFarmTokenAttributesType;
};

export type UnstakeEventType = {
    farmingToken: DcdtTokenPaymentType;
    farmToken: DcdtTokenPaymentType;
    farmSupply: string;
    rewardTokens: DcdtTokenPaymentType;
    rewardTokenReserves: string;
    farmAttributes: StakingFarmTokenAttributesType;
};

export type StakeClaimRewardsEventType = {
    oldFarmToken: DcdtTokenPaymentType;
    newFarmToken: DcdtTokenPaymentType;
    farmSupply: string;
    rewardTokens: DcdtTokenPaymentType;
    rewardTokenReserves: string;
    oldFarmAttributes: StakingFarmTokenAttributesType;
    newFarmAttributes: StakingFarmTokenAttributesType;
    createdWithMerge: boolean;
};
