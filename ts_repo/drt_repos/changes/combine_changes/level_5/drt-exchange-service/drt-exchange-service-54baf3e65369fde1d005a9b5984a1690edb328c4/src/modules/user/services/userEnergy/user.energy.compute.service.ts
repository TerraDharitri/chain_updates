import { Injectable } from '@nestjs/common';
import { scAddress } from '../../../../config';
import {
    EnergyType,
    LockedAssetAttributes,
    LockedTokenAttributes,
    WrappedFarmTokenAttributes,
    WrappedFarmTokenAttributesV2,
    WrappedLpTokenAttributes,
    WrappedLpTokenAttributesV2,
} from '@terradharitri/sdk-exchange';
import { ClaimProgress } from '../../../../submodules/weekly-rewards-splitting/models/weekly-rewards-splitting.model';
import {
    ContractType,
    OutdatedContract,
    UserNegativeEnergyCheck,
} from '../../models/user.model';
import { ProxyService } from '../../../proxy/services/proxy.service';
import { StakingProxyService } from '../../../staking-proxy/services/staking.proxy.service';
import { FarmVersion } from '../../../farm/models/farm.model';
import { farmVersion } from '../../../../utils/farm.utils';
import { BigNumber } from 'bignumber.js';
import { WeekTimekeepingAbiService } from 'src/submodules/week-timekeeping/services/week-timekeeping.abi.service';
import { WeeklyRewardsSplittingAbiService } from 'src/submodules/weekly-rewards-splitting/services/weekly-rewards-splitting.abi.service';
import { StakingProxyAbiService } from 'src/modules/staking-proxy/services/staking.proxy.abi.service';
import { FarmAbiFactory } from 'src/modules/farm/farm.abi.factory';
import { FarmFactoryService } from 'src/modules/farm/farm.factory';
import { FarmServiceV2 } from 'src/modules/farm/v2/services/farm.v2.service';
import { GetOrSetCache } from 'src/helpers/decorators/caching.decorator';
import { Constants, ErrorLoggerAsync } from '@terradharitri/sdk-nestjs-common';
import { EnergyAbiService } from 'src/modules/energy/services/energy.abi.service';
import { RemoteConfigGetterService } from 'src/modules/remote-config/remote-config.getter.service';
import { StakingService } from 'src/modules/staking/services/staking.service';
import { MXApiService } from 'src/services/dharitri-communication/drt.api.service';
import { tokenIdentifier } from 'src/utils/token.converters';
import { ProxyPairAbiService } from 'src/modules/proxy/services/proxy-pair/proxy.pair.abi.service';
import { ProxyFarmAbiService } from 'src/modules/proxy/services/proxy-farm/proxy.farm.abi.service';
import { LockedAssetGetterService } from 'src/modules/locked-asset-factory/services/locked.asset.getter.service';
import { MetabondingAbiService } from 'src/modules/metabonding/services/metabonding.abi.service';
import { PaginationArgs } from 'src/modules/dex.model';
import { UserMetaDcdtService } from '../user.metaDcdt.service';
import { StakingComputeService } from 'src/modules/staking/services/staking.compute.service';

@Injectable()
export class UserEnergyComputeService {
    constructor(
        private readonly farmAbi: FarmAbiFactory,
        private readonly farmService: FarmFactoryService,
        private readonly weekTimekeepingAbi: WeekTimekeepingAbiService,
        private readonly weeklyRewardsSplittingAbi: WeeklyRewardsSplittingAbiService,
        private readonly stakeProxyService: StakingProxyService,
        private readonly stakeProxyAbi: StakingProxyAbiService,
        private readonly stakingService: StakingService,
        private readonly stakingCompute: StakingComputeService,
        private readonly energyAbi: EnergyAbiService,
        private readonly lockedAssetGetter: LockedAssetGetterService,
        private readonly proxyService: ProxyService,
        private readonly proxyPairAbi: ProxyPairAbiService,
        private readonly proxyFarmAbi: ProxyFarmAbiService,
        private readonly metabondingAbi: MetabondingAbiService,
        private readonly drtApi: MXApiService,
        private readonly remoteConfig: RemoteConfigGetterService,
        private readonly userMetaDcdtService: UserMetaDcdtService,
    ) {}

    async getUserOutdatedContracts(
        userAddress: string,
        skipFeesCollector = false,
    ): Promise<OutdatedContract[]> {
        const activeFarms = await this.userActiveFarmsV2(userAddress);

        const promises = activeFarms.map((farm) =>
            this.outdatedContract(userAddress, farm),
        );

        const activeStakings = await this.userActiveStakings(userAddress);
        promises.push(
            ...activeStakings.map((stake) =>
                this.outdatedContract(userAddress, stake),
            ),
        );

        if (!skipFeesCollector) {
            promises.push(
                this.outdatedContract(userAddress, scAddress.feesCollector),
            );
        }

        const outdatedContracts = await Promise.all(promises);
        return outdatedContracts.filter(
            (contract) => contract && contract.address,
        );
    }

    async outdatedContract(
        userAddress: string,
        contractAddress: string,
    ): Promise<OutdatedContract> {
        return await this.computeUserOutdatedContract(
            userAddress,
            contractAddress,
        );
    }

    async computeUserOutdatedContract(
        userAddress: string,
        contractAddress: string,
    ): Promise<OutdatedContract> {
        const stakeAddresses = await this.remoteConfig.getStakingAddresses();
        if (stakeAddresses.includes(contractAddress)) {
            return await this.computeStakingOutdatedContract(
                userAddress,
                contractAddress,
            );
        }

        const isFarmAddress = contractAddress !== scAddress.feesCollector;

        if (isFarmAddress) {
            return await this.computeFarmOutdatedContract(
                userAddress,
                contractAddress,
            );
        }

        return await this.computeFeesCollectorOutdatedContract(userAddress);
    }

    async computeFarmOutdatedContract(
        userAddress: string,
        contractAddress: string,
    ): Promise<OutdatedContract> {
        const farmService = this.farmService.useService(
            contractAddress,
        ) as FarmServiceV2;
        const [
            currentClaimProgress,
            currentWeek,
            farmToken,
            userEnergy,
            produceRewardsEnabled,
        ] = await Promise.all([
            this.weeklyRewardsSplittingAbi.currentClaimProgress(
                contractAddress,
                userAddress,
            ),
            this.weekTimekeepingAbi.currentWeek(contractAddress),
            farmService.getFarmToken(contractAddress),
            this.energyAbi.energyEntryForUser(userAddress),
            this.farmAbi
                .useAbi(contractAddress)
                .produceRewardsEnabled(contractAddress),
        ]);

        if (currentClaimProgress.week === 0) {
            return new OutdatedContract();
        }

        const claimProgressTotalRewards =
            await this.weeklyRewardsSplittingAbi.totalRewardsForWeek(
                contractAddress,
                currentClaimProgress.week,
            );

        const outdatedClaimProgress = currentClaimProgress.week !== currentWeek;

        if (
            !produceRewardsEnabled &&
            outdatedClaimProgress &&
            (claimProgressTotalRewards.length === 0 ||
                claimProgressTotalRewards[0].amount === '0')
        ) {
            return new OutdatedContract();
        }

        if (
            this.isEnergyOutdated(userEnergy, currentClaimProgress) ||
            outdatedClaimProgress
        ) {
            return new OutdatedContract({
                address: contractAddress,
                type: ContractType.Farm,
                claimProgressOutdated: outdatedClaimProgress,
                farmToken: farmToken.collection,
            });
        }

        return new OutdatedContract();
    }

    async computeStakingOutdatedContract(
        userAddress: string,
        contractAddress: string,
    ): Promise<OutdatedContract> {
        const [
            currentClaimProgress,
            currentWeek,
            farmToken,
            userEnergy,
            isProducingRewards,
        ] = await Promise.all([
            this.weeklyRewardsSplittingAbi.currentClaimProgress(
                contractAddress,
                userAddress,
            ),
            this.weekTimekeepingAbi.currentWeek(contractAddress),
            this.stakingService.getFarmToken(contractAddress),
            this.energyAbi.energyEntryForUser(userAddress),
            this.stakingCompute.isProducingRewards(contractAddress),
        ]);

        if (currentClaimProgress.week === 0) {
            return new OutdatedContract();
        }

        const claimProgressTotalRewards =
            await this.weeklyRewardsSplittingAbi.totalRewardsForWeek(
                contractAddress,
                currentClaimProgress.week,
            );

        const outdatedClaimProgress = currentClaimProgress.week !== currentWeek;

        if (
            !isProducingRewards &&
            outdatedClaimProgress &&
            (claimProgressTotalRewards.length === 0 ||
                claimProgressTotalRewards[0].amount === '0')
        ) {
            return new OutdatedContract();
        }

        if (
            this.isEnergyOutdated(userEnergy, currentClaimProgress) ||
            outdatedClaimProgress
        ) {
            return new OutdatedContract({
                address: contractAddress,
                type: ContractType.StakingFarm,
                claimProgressOutdated: outdatedClaimProgress,
                farmToken: farmToken.collection,
            });
        }

        return new OutdatedContract();
    }

    async computeFeesCollectorOutdatedContract(
        userAddress: string,
    ): Promise<OutdatedContract> {
        const [currentClaimProgress, currentWeek, userEnergy] =
            await Promise.all([
                this.weeklyRewardsSplittingAbi.currentClaimProgress(
                    scAddress.feesCollector,
                    userAddress,
                ),
                this.weekTimekeepingAbi.currentWeek(scAddress.feesCollector),
                this.energyAbi.energyEntryForUser(userAddress),
            ]);

        if (currentClaimProgress.week === 0) {
            return new OutdatedContract();
        }

        const outdatedClaimProgress = currentClaimProgress.week !== currentWeek;

        if (
            this.isEnergyOutdated(userEnergy, currentClaimProgress) ||
            outdatedClaimProgress
        ) {
            return new OutdatedContract({
                address: scAddress.feesCollector,
                type: ContractType.FeesCollector,
                claimProgressOutdated: outdatedClaimProgress,
            });
        }
        return new OutdatedContract();
    }

    @GetOrSetCache({
        baseKey: 'userEnergy',
        remoteTtl: Constants.oneMinute(),
    })
    async userActiveFarmsV2(userAddress: string): Promise<string[]> {
        return await this.computeActiveFarmsV2ForUser(userAddress);
    }

    async computeActiveFarmsV2ForUser(userAddress: string): Promise<string[]> {
        const maxPagination = new PaginationArgs({
            limit: 100,
            offset: 0,
        });
        const [farmTokens, farmLockedTokens, dualYieldTokens] =
            await Promise.all([
                this.userMetaDcdtService.getUserFarmTokens(
                    userAddress,
                    maxPagination,
                    false,
                ),
                this.userMetaDcdtService.getUserLockedFarmTokensV2(
                    userAddress,
                    maxPagination,
                    false,
                ),
                this.userMetaDcdtService.getUserDualYieldTokens(
                    userAddress,
                    maxPagination,
                    false,
                ),
            ]);

        let userActiveFarmAddresses = farmTokens.map((token) => token.creator);
        const promisesFarmLockedTokens = farmLockedTokens.map((token) => {
            return this.decodeAndGetFarmAddressFarmLockedTokens(
                token.identifier,
                token.attributes,
            );
        });
        const promisesDualYieldTokens = dualYieldTokens.map((token) => {
            return this.getFarmAddressForDualYieldToken(token.collection);
        });

        userActiveFarmAddresses = userActiveFarmAddresses.concat(
            await Promise.all([
                ...promisesFarmLockedTokens,
                ...promisesDualYieldTokens,
            ]),
        );
        return [...new Set(userActiveFarmAddresses)].filter(
            (address) => farmVersion(address) === FarmVersion.V2,
        );
    }

    @GetOrSetCache({
        baseKey: 'userEnergy',
        remoteTtl: Constants.oneMinute(),
    })
    async userActiveStakings(userAddress: string): Promise<string[]> {
        return await this.computeActiveStakingsForUser(userAddress);
    }

    async computeActiveStakingsForUser(userAddress: string): Promise<string[]> {
        const maxPagination = new PaginationArgs({
            limit: 100,
            offset: 0,
        });
        const [stakeTokens, dualYieldTokens] = await Promise.all([
            this.userMetaDcdtService.getUserStakeFarmTokens(
                userAddress,
                maxPagination,
            ),
            this.userMetaDcdtService.getUserDualYieldTokens(
                userAddress,
                maxPagination,
                false,
            ),
        ]);

        let userActiveStakeAddresses = stakeTokens.map(
            (token) => token.creator,
        );
        const promisesDualYieldTokens = dualYieldTokens.map((token) => {
            return this.getStakeAddressForDualYieldToken(token.collection);
        });

        userActiveStakeAddresses = userActiveStakeAddresses.concat(
            await Promise.all([...promisesDualYieldTokens]),
        );
        return [...new Set(userActiveStakeAddresses)];
    }

    @ErrorLoggerAsync({
        logArgs: true,
    })
    async computeNegativeEnergyCheck(
        userAddress: string,
    ): Promise<UserNegativeEnergyCheck> {
        const userNftsCount = await this.drtApi.getNftsCountForUser(userAddress);

        if (userNftsCount === 0) {
            return new UserNegativeEnergyCheck({
                LKMOA: false,
                XMOA: false,
                lockedLPTokenV1: false,
                lockedLPTokenV2: false,
                lockedFarmTokenV1: false,
                lockedFarmTokenV2: false,
                metabonding: false,
            });
        }

        const [
            lkmoaTokenID,
            xmoaTokenID,
            lockedLPTokenIDV1,
            lockedLPTokenIDV2,
            lockedFarmTokenIDV1,
            lockedFarmTokenIDV2,
            lkmoaActivationNonce,
            stats,
        ] = await Promise.all([
            this.lockedAssetGetter.getLockedTokenID(),
            this.energyAbi.lockedTokenID(),
            this.proxyPairAbi.wrappedLpTokenID(scAddress.proxyDexAddress.v1),
            this.proxyPairAbi.wrappedLpTokenID(scAddress.proxyDexAddress.v2),
            this.proxyFarmAbi.wrappedFarmTokenID(scAddress.proxyDexAddress.v1),
            this.proxyFarmAbi.wrappedFarmTokenID(scAddress.proxyDexAddress.v2),
            this.lockedAssetGetter.getExtendedAttributesActivationNonce(),
            this.drtApi.getStats(),
        ]);

        const userNfts = await this.drtApi.getNftsForUser(
            userAddress,
            'MetaDCDT',
        );

        const lkmoaTokens = userNfts.filter((nft) =>
            nft.collection.includes(lkmoaTokenID),
        );
        const xmoaTokens = userNfts.filter((nft) =>
            nft.collection.includes(xmoaTokenID),
        );
        const lockedLPTokensV1 = userNfts.filter(
            (nft) => nft.collection === lockedLPTokenIDV1,
        );
        const lockedLPTokensV2 = userNfts.filter(
            (nft) => nft.collection === lockedLPTokenIDV2,
        );
        const lockedFarmTokensV1 = userNfts.filter(
            (nft) => nft.collection === lockedFarmTokenIDV1,
        );
        const lockedFarmTokensV2 = userNfts.filter(
            (nft) => nft.collection === lockedFarmTokenIDV2,
        );

        const [
            checkLKMOA,
            checkXMOA,
            checkLockedLPTokenV1,
            checkLockedLPTokenV2,
            checkLockedFarmTokenV1,
            checkLockedFarmTokenV2,
        ] = await Promise.all([
            this.checkLKMOANegativeEnergy(
                stats.epoch,
                lkmoaTokens.map((token) =>
                    LockedAssetAttributes.fromAttributes(
                        token.nonce >= lkmoaActivationNonce,
                        token.attributes,
                    ),
                ),
            ),
            this.checkXMOANegativeEnergy(
                stats.epoch,
                xmoaTokens.map((token) =>
                    LockedTokenAttributes.fromAttributes(token.attributes),
                ),
            ),
            this.checkLockedLPTokenNegativeEnergyV1(
                stats.epoch,
                lkmoaActivationNonce,
                lkmoaTokenID,
                lockedLPTokensV1.map((token) =>
                    WrappedLpTokenAttributes.fromAttributes(token.attributes),
                ),
            ),
            this.checkLockedLPTokenNegativeEnergyV2(
                stats.epoch,
                lkmoaActivationNonce,
                lkmoaTokenID,
                lockedLPTokensV2.map((token) =>
                    WrappedLpTokenAttributesV2.fromAttributes(token.attributes),
                ),
            ),
            this.checkLockedFarmTokenNegativeEnergyV1(
                stats.epoch,
                lkmoaActivationNonce,
                lkmoaTokenID,
                lockedFarmTokensV1.map((token) =>
                    WrappedFarmTokenAttributes.fromAttributes(token.attributes),
                ),
            ),
            this.checkLockedFarmTokenNegativeEnergyV2(
                stats.epoch,
                lkmoaActivationNonce,
                lkmoaTokenID,
                lockedFarmTokensV2.map((token) =>
                    WrappedFarmTokenAttributesV2.fromAttributes(
                        token.attributes,
                    ),
                ),
            ),
        ]);

        let metabondingCheck = false;
        if (scAddress.metabondingStakingAddress) {
            const userMetabondingEntry = await this.metabondingAbi.userEntry(
                userAddress,
            );

            if (userMetabondingEntry.tokenNonce > 0) {
                const metabondingTokensAttributes =
                    await this.drtApi.getNftAttributesByTokenIdentifier(
                        scAddress.metabondingStakingAddress,
                        tokenIdentifier(
                            lkmoaTokenID,
                            userMetabondingEntry.tokenNonce,
                        ),
                    );
                metabondingCheck = this.checkLKMOANegativeEnergy(stats.epoch, [
                    LockedAssetAttributes.fromAttributes(
                        userMetabondingEntry.tokenNonce >= lkmoaActivationNonce,
                        metabondingTokensAttributes,
                    ),
                ]);
            }
        }

        return new UserNegativeEnergyCheck({
            LKMOA: checkLKMOA,
            XMOA: checkXMOA,
            lockedLPTokenV1: checkLockedLPTokenV1,
            lockedLPTokenV2: checkLockedLPTokenV2,
            lockedFarmTokenV1: checkLockedFarmTokenV1,
            lockedFarmTokenV2: checkLockedFarmTokenV2,
            metabonding: metabondingCheck,
        });
    }

    checkLKMOANegativeEnergy(
        currentEpoch: number,
        lockedTokensAttributes: LockedAssetAttributes[],
    ): boolean {
        for (const attributes of lockedTokensAttributes) {
            for (const schedule of attributes.unlockSchedule) {
                if (schedule.epoch.integerValue().toNumber() > currentEpoch) {
                    return false;
                }
            }
        }
        if (lockedTokensAttributes.length === 0) {
            return false;
        }
        return true;
    }

    checkXMOANegativeEnergy(
        currentEpoch: number,
        lockedTokensAttributes: LockedTokenAttributes[],
    ): boolean {
        for (const attributes of lockedTokensAttributes) {
            if (attributes.unlockEpoch > currentEpoch) {
                return false;
            }
        }
        if (lockedTokensAttributes.length === 0) {
            return false;
        }
        return true;
    }

    async checkLockedLPTokenNegativeEnergyV1(
        currentEpoch: number,
        lkmoaExtendedAttributesActivationNonce: number,
        lkmoaTokenID: string,
        attributes: WrappedLpTokenAttributes[],
    ): Promise<boolean> {
        const identifiers = attributes.map((attribute) =>
            tokenIdentifier(lkmoaTokenID, attribute.lockedAssetsNonce),
        );
        const lockedTokensAttributes =
            await this.drtApi.getNftsAttributesForUser(
                scAddress.proxyDexAddress.v1,
                'MetaDCDT',
                identifiers,
            );

        return this.checkLKMOANegativeEnergy(
            currentEpoch,
            lockedTokensAttributes.map((attribute, index) =>
                LockedAssetAttributes.fromAttributes(
                    attributes[index].lockedAssetsNonce >=
                        lkmoaExtendedAttributesActivationNonce,
                    attribute,
                ),
            ),
        );
    }

    async checkLockedLPTokenNegativeEnergyV2(
        currentEpoch: number,
        lkmoaExtendedAttributesActivationNonce: number,
        lkmoaTokenID: string,
        attributes: WrappedLpTokenAttributesV2[],
    ): Promise<boolean> {
        const lockedLPWithLKMOA = attributes.filter(
            (attribute) =>
                attribute.lockedTokens.tokenIdentifier === lkmoaTokenID,
        );
        const lockedLPWithXMOA = attributes.filter(
            (attribute) =>
                attribute.lockedTokens.tokenIdentifier !== lkmoaTokenID,
        );

        const [lkmoaAttributesRaw, xmoaAttributesRaw] = await Promise.all([
            this.drtApi.getNftsAttributesForUser(
                scAddress.proxyDexAddress.v2,
                'MetaDCDT',
                lockedLPWithLKMOA.map((attribute) =>
                    tokenIdentifier(
                        attribute.lockedTokens.tokenIdentifier,
                        attribute.lockedTokens.tokenNonce,
                    ),
                ),
            ),
            this.drtApi.getNftsAttributesForUser(
                scAddress.proxyDexAddress.v2,
                'MetaDCDT',
                lockedLPWithXMOA.map((attribute) =>
                    tokenIdentifier(
                        attribute.lockedTokens.tokenIdentifier,
                        attribute.lockedTokens.tokenNonce,
                    ),
                ),
            ),
        ]);

        const checkLKMOA = this.checkLKMOANegativeEnergy(
            currentEpoch,
            lkmoaAttributesRaw.map((attribute, index) =>
                LockedAssetAttributes.fromAttributes(
                    lockedLPWithLKMOA[index].lockedTokens.tokenNonce >=
                        lkmoaExtendedAttributesActivationNonce,
                    attribute,
                ),
            ),
        );
        const checkXMOA = this.checkXMOANegativeEnergy(
            currentEpoch,
            xmoaAttributesRaw.map((attribute) =>
                LockedTokenAttributes.fromAttributes(attribute),
            ),
        );

        return checkLKMOA || checkXMOA;
    }

    async checkLockedFarmTokenNegativeEnergyV1(
        currentEpoch: number,
        lkmoaExtendedAttributesActivationNonce: number,
        lkmoaTokenID: string,
        attributes: WrappedFarmTokenAttributes[],
    ): Promise<boolean> {
        const lockedFarmTokensWithLKMOA = attributes.filter(
            (attribute) => attribute.farmingTokenID === lkmoaTokenID,
        );
        const lockedFarmTokensWithLKLP = attributes.filter(
            (attribute) => attribute.farmingTokenID !== lkmoaTokenID,
        );

        const [lkmoaAttributesRaw, lklpAttributesRaw] = await Promise.all([
            this.drtApi.getNftsAttributesForUser(
                scAddress.proxyDexAddress.v1,
                'MetaDCDT',
                lockedFarmTokensWithLKMOA.map((attribute) =>
                    tokenIdentifier(
                        attribute.farmingTokenID,
                        attribute.farmingTokenNonce,
                    ),
                ),
            ),
            this.drtApi.getNftsAttributesForUser(
                scAddress.proxyDexAddress.v1,
                'MetaDCDT',
                lockedFarmTokensWithLKLP.map((attribute) =>
                    tokenIdentifier(
                        attribute.farmingTokenID,
                        attribute.farmingTokenNonce,
                    ),
                ),
            ),
        ]);

        const checkLKMOA = this.checkLKMOANegativeEnergy(
            currentEpoch,
            lkmoaAttributesRaw.map((attribute, index) => {
                return LockedAssetAttributes.fromAttributes(
                    lockedFarmTokensWithLKMOA[index].farmingTokenNonce >=
                        lkmoaExtendedAttributesActivationNonce,
                    attribute,
                );
            }),
        );

        const checkLockedLPTokenV1 =
            await this.checkLockedLPTokenNegativeEnergyV1(
                currentEpoch,
                lkmoaExtendedAttributesActivationNonce,
                lkmoaTokenID,
                lklpAttributesRaw.map((attribute) =>
                    WrappedLpTokenAttributes.fromAttributes(attribute),
                ),
            );

        return checkLKMOA || checkLockedLPTokenV1;
    }

    async checkLockedFarmTokenNegativeEnergyV2(
        currentEpoch: number,
        lkmoaExtendedAttributesActivationNonce: number,
        lkmoaTokenID: string,
        attributes: WrappedFarmTokenAttributesV2[],
    ): Promise<boolean> {
        const lklpAttributesRaw = await this.drtApi.getNftsAttributesForUser(
            scAddress.proxyDexAddress.v2,
            'MetaDCDT',
            attributes.map((attribute) =>
                tokenIdentifier(
                    attribute.proxyFarmingToken.tokenIdentifier,
                    attribute.proxyFarmingToken.tokenNonce,
                ),
            ),
        );

        return this.checkLockedLPTokenNegativeEnergyV2(
            currentEpoch,
            lkmoaExtendedAttributesActivationNonce,
            lkmoaTokenID,
            lklpAttributesRaw.map((attribute) =>
                WrappedLpTokenAttributesV2.fromAttributes(attribute),
            ),
        );
    }

    decodeAndGetFarmAddressFarmLockedTokens(
        identifier: string,
        attributes: string,
    ): Promise<string> {
        const decodedWFMTAttributes =
            this.proxyService.getWrappedFarmTokenAttributesV2({
                batchAttributes: [
                    {
                        identifier,
                        attributes,
                    },
                ],
            });

        return this.farmAbi.getFarmAddressByFarmTokenID(
            decodedWFMTAttributes[0].farmToken.tokenIdentifier,
        );
    }

    async getFarmAddressForDualYieldToken(collection: string): Promise<string> {
        if (!collection || collection === undefined) {
            return undefined;
        }

        const stakingProxyAddress =
            await this.stakeProxyService.getStakingProxyAddressByDualYieldTokenID(
                collection,
            );
        return this.stakeProxyAbi.lpFarmAddress(stakingProxyAddress);
    }

    async getStakeAddressForDualYieldToken(collection: string) {
        if (!collection || collection === undefined) {
            return undefined;
        }

        const stakingProxyAddress =
            await this.stakeProxyService.getStakingProxyAddressByDualYieldTokenID(
                collection,
            );
        return this.stakeProxyAbi.stakingFarmAddress(stakingProxyAddress);
    }

    isEnergyOutdated(
        currentUserEnergy: EnergyType,
        currentClaimProgress: ClaimProgress,
    ): boolean {
        if (currentClaimProgress.week === 0) {
            return false;
        }

        if (
            currentUserEnergy.lastUpdateEpoch >
            currentClaimProgress.energy.lastUpdateEpoch
        ) {
            const epochsDiff =
                currentUserEnergy.lastUpdateEpoch -
                currentClaimProgress.energy.lastUpdateEpoch;
            currentClaimProgress.energy.amount = new BigNumber(
                currentClaimProgress.energy.amount,
            )
                .minus(
                    new BigNumber(epochsDiff).multipliedBy(
                        currentClaimProgress.energy.totalLockedTokens,
                    ),
                )
                .toFixed();
        }

        if (
            currentClaimProgress.energy.lastUpdateEpoch >
            currentUserEnergy.lastUpdateEpoch
        ) {
            const epochsDiff =
                currentClaimProgress.energy.lastUpdateEpoch -
                currentUserEnergy.lastUpdateEpoch;
            currentUserEnergy.amount = new BigNumber(currentUserEnergy.amount)
                .minus(
                    new BigNumber(epochsDiff).multipliedBy(
                        currentUserEnergy.totalLockedTokens,
                    ),
                )
                .toFixed();
        }
        return currentUserEnergy.amount !== currentClaimProgress.energy.amount;
    }
}
