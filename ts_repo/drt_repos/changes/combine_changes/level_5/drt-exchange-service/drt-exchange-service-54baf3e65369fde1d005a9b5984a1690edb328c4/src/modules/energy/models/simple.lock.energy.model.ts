import { Field, Int, ObjectType } from '@nestjs/graphql';
import { nestedFieldComplexity } from 'src/helpers/complexity/field.estimators';
import { DcdtToken } from 'src/modules/tokens/models/dcdtToken.model';
import { NftCollection } from 'src/modules/tokens/models/nftCollection.model';

@ObjectType()
export class LockOption {
    @Field(() => Int)
    lockEpochs: number;
    @Field(() => Int)
    penaltyStartPercentage: number;

    constructor(init: LockOption) {
        Object.assign(this, init);
    }
}

@ObjectType()
export class SimpleLockEnergyModel {
    @Field()
    address: string;
    @Field({ complexity: nestedFieldComplexity })
    baseAssetToken: DcdtToken;
    @Field({ complexity: nestedFieldComplexity })
    lockedToken: NftCollection;
    @Field({ complexity: nestedFieldComplexity })
    legacyLockedToken: NftCollection;
    @Field(() => [LockOption], { complexity: nestedFieldComplexity })
    lockOptions: LockOption[];
    @Field()
    tokenUnstakeAddress: string;
    @Field()
    pauseState: boolean;

    constructor(init?: Partial<SimpleLockEnergyModel>) {
        Object.assign(this, init);
    }
}
