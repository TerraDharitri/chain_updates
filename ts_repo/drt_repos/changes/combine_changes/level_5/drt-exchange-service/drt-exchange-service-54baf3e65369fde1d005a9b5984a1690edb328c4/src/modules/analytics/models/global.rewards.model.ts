import { Field, ObjectType } from '@nestjs/graphql';
import { nestedFieldComplexity } from 'src/helpers/complexity/field.estimators';
import { DcdtToken } from 'src/modules/tokens/models/dcdtToken.model';

@ObjectType()
export class FeesCollectorGlobalRewards {
    @Field()
    totalRewardsUSD: string;

    @Field()
    energyRewardsUSD: string;

    constructor(init?: Partial<FeesCollectorGlobalRewards>) {
        Object.assign(this, init);
    }
}

@ObjectType()
export class FarmsGlobalRewards {
    @Field()
    pairAddress: string;

    @Field(() => DcdtToken, { nullable: true, complexity: nestedFieldComplexity })
    firstToken?: DcdtToken;

    @Field(() => DcdtToken, { nullable: true, complexity: nestedFieldComplexity })
    secondToken?: DcdtToken;

    @Field()
    totalRewardsUSD: string;

    @Field()
    energyRewardsUSD: string;

    constructor(init?: Partial<FarmsGlobalRewards>) {
        Object.assign(this, init);
        if (init?.firstToken) {
            this.firstToken = new DcdtToken(init.firstToken);
        }
        if (init?.secondToken) {
            this.secondToken = new DcdtToken(init.secondToken);
        }
    }
}

@ObjectType()
export class StakingGlobalRewards {
    @Field(() => DcdtToken, { nullable: true, complexity: nestedFieldComplexity })
    farmingToken?: DcdtToken;

    @Field()
    totalRewardsUSD: string;

    @Field()
    energyRewardsUSD: string;

    constructor(init?: Partial<StakingGlobalRewards>) {
        Object.assign(this, init);
        if (init?.farmingToken) {
            this.farmingToken = new DcdtToken(init.farmingToken);
        }
    }
}

@ObjectType()
export class GlobalRewardsModel {
    @Field(() => FeesCollectorGlobalRewards)
    feesCollectorGlobalRewards: FeesCollectorGlobalRewards;

    @Field(() => [FarmsGlobalRewards])
    farmsGlobalRewards: FarmsGlobalRewards[];

    @Field(() => [StakingGlobalRewards])
    stakingGlobalRewards: StakingGlobalRewards[];

    constructor(init?: Partial<GlobalRewardsModel>) {
        Object.assign(this, init);
        if (init?.feesCollectorGlobalRewards) {
            this.feesCollectorGlobalRewards = new FeesCollectorGlobalRewards(
                init.feesCollectorGlobalRewards,
            );
        }
        if (init?.farmsGlobalRewards) {
            this.farmsGlobalRewards = init.farmsGlobalRewards.map(
                (farm) => new FarmsGlobalRewards(farm),
            );
        }
        if (init?.stakingGlobalRewards) {
            this.stakingGlobalRewards = init.stakingGlobalRewards.map(
                (staking) => new StakingGlobalRewards(staking),
            );
        }
    }
}
