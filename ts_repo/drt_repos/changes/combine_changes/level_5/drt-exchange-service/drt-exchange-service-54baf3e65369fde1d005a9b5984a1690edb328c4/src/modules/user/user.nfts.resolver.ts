import { UseGuards } from '@nestjs/common';
import { Args, Query, ResolveField, Resolver } from '@nestjs/graphql';
import { AuthUser } from '../auth/auth.user';
import { UserAuthResult } from '../auth/user.auth.result';
import { JwtOrNativeAuthGuard } from '../auth/jwt.or.native.auth.guard';
import { PaginationArgs } from '../dex.model';
import {
    UserDualYiledToken,
    UserFarmToken,
    UserLockedAssetToken,
    UserLockedDcdtToken,
    UserLockedFarmToken,
    UserLockedFarmTokenV2,
    UserLockedLPToken,
    UserLockedLPTokenV2,
    UserLockedSimpleFarmToken,
    UserLockedSimpleLpToken,
    UserLockedTokenEnergy,
    UserNftsModel,
    UserRedeemToken,
    UserStakeFarmToken,
    UserUnbondFarmToken,
    UserWrappedLockedToken,
} from './models/user.model';
import { UserMetaDcdtService } from './services/user.metaDcdt.service';
import { ContextGetterService } from 'src/services/context/context.getter.service';

@Resolver(() => UserNftsModel)
export class UserNftsResolver {
    constructor(
        private readonly userMetaDcdts: UserMetaDcdtService,
        private readonly contextGetter: ContextGetterService,
    ) {}

    @ResolveField()
    async userLockedAssetToken(
        parent: UserNftsModel,
    ): Promise<UserLockedAssetToken[]> {
        return this.userMetaDcdts.getUserLockedAssetTokens(
            parent.address,
            parent.pagination,
            parent.rawNfts,
        );
    }

    @ResolveField()
    async userFarmToken(parent: UserNftsModel): Promise<UserFarmToken[]> {
        return this.userMetaDcdts.getUserFarmTokens(
            parent.address,
            parent.pagination,
            true,
            parent.rawNfts,
        );
    }

    @ResolveField()
    async userLockedLPToken(
        parent: UserNftsModel,
    ): Promise<UserLockedLPToken[]> {
        return this.userMetaDcdts.getUserLockedLpTokens(
            parent.address,
            parent.pagination,
            parent.rawNfts,
        );
    }

    @ResolveField()
    async userLockedFarmToken(
        parent: UserNftsModel,
    ): Promise<UserLockedFarmToken[]> {
        return this.userMetaDcdts.getUserLockedFarmTokens(
            parent.address,
            parent.pagination,
            parent.rawNfts,
        );
    }

    @ResolveField()
    async userLockedLpTokenV2(
        parent: UserNftsModel,
    ): Promise<UserLockedLPTokenV2[]> {
        return this.userMetaDcdts.getUserLockedLpTokensV2(
            parent.address,
            parent.pagination,
            parent.rawNfts,
        );
    }

    @ResolveField()
    async userLockedFarmTokenV2(
        parent: UserNftsModel,
    ): Promise<UserLockedFarmTokenV2[]> {
        return this.userMetaDcdts.getUserLockedFarmTokensV2(
            parent.address,
            parent.pagination,
            true,
            parent.rawNfts,
        );
    }

    @ResolveField()
    async userStakeFarmToken(
        parent: UserNftsModel,
    ): Promise<UserStakeFarmToken[]> {
        return this.userMetaDcdts.getUserStakeFarmTokens(
            parent.address,
            parent.pagination,
            parent.rawNfts,
        );
    }

    @ResolveField()
    async userUnbondFarmToken(
        parent: UserNftsModel,
    ): Promise<UserUnbondFarmToken[]> {
        return this.userMetaDcdts.getUserUnbondFarmTokens(
            parent.address,
            parent.pagination,
            parent.rawNfts,
        );
    }

    @ResolveField()
    async userDualYieldToken(
        parent: UserNftsModel,
    ): Promise<UserDualYiledToken[]> {
        return this.userMetaDcdts.getUserDualYieldTokens(
            parent.address,
            parent.pagination,
            true,
            parent.rawNfts,
        );
    }

    @ResolveField()
    async userRedeemToken(parent: UserNftsModel): Promise<UserRedeemToken[]> {
        return this.userMetaDcdts.getUserRedeemToken(
            parent.address,
            parent.pagination,
            parent.rawNfts,
        );
    }

    @ResolveField()
    async userLockedDcdtToken(
        parent: UserNftsModel,
    ): Promise<UserLockedDcdtToken[]> {
        return this.userMetaDcdts.getUserLockedDcdtToken(
            parent.address,
            parent.pagination,
            parent.rawNfts,
        );
    }

    @ResolveField()
    async userLockedSimpleLpToken(
        parent: UserNftsModel,
    ): Promise<UserLockedSimpleLpToken[]> {
        return this.userMetaDcdts.getUserLockedSimpleLpToken(
            parent.address,
            parent.pagination,
            parent.rawNfts,
        );
    }

    @ResolveField()
    async userLockedSimpleFarmToken(
        parent: UserNftsModel,
    ): Promise<UserLockedSimpleFarmToken[]> {
        return this.userMetaDcdts.getUserLockedSimpleFarmToken(
            parent.address,
            parent.pagination,
            parent.rawNfts,
        );
    }

    @ResolveField()
    async userLockedTokenEnergy(
        parent: UserNftsModel,
    ): Promise<UserLockedTokenEnergy[]> {
        return this.userMetaDcdts.getUserLockedTokenEnergy(
            parent.address,
            parent.pagination,
            parent.rawNfts,
        );
    }

    @ResolveField()
    async userWrappedLockedToken(
        parent: UserNftsModel,
    ): Promise<UserWrappedLockedToken[]> {
        return this.userMetaDcdts.getUserWrappedLockedTokenEnergy(
            parent.address,
            parent.pagination,
            parent.rawNfts,
        );
    }

    @UseGuards(JwtOrNativeAuthGuard)
    @Query(() => UserNftsModel)
    async userNfts(
        @Args() pagination: PaginationArgs,
        @AuthUser() user: UserAuthResult,
    ): Promise<UserNftsModel> {
        const nfts = await this.contextGetter.getNftsForUser(
            user.address,
            pagination.offset,
            pagination.limit,
            'MetaDCDT',
        );

        return new UserNftsModel(user.address, pagination, nfts);
    }
}
