import { Field, Int, ObjectType } from '@nestjs/graphql';
import { DcdtTokenPaymentModel } from 'src/modules/tokens/models/dcdt.token.payment.model';

@ObjectType()
export class TokenUnstakeModel {
    @Field()
    address: string;
    @Field(() => Int)
    unbondEpochs: number;
    @Field(() => Int)
    feesBurnPercentage: number;
    @Field()
    feesCollectorAddress: string;
    @Field()
    energyFactoryAddress: string;

    constructor(init?: Partial<TokenUnstakeModel>) {
        Object.assign(this, init);
    }
}

@ObjectType()
export class UnstakePairModel {
    @Field()
    unlockEpoch: number;
    @Field(() => DcdtTokenPaymentModel)
    lockedTokens: DcdtTokenPaymentModel;
    @Field(() => DcdtTokenPaymentModel)
    unlockedTokens: DcdtTokenPaymentModel;

    constructor(init?: Partial<UnstakePairModel>) {
        Object.assign(this, init);
    }
}
