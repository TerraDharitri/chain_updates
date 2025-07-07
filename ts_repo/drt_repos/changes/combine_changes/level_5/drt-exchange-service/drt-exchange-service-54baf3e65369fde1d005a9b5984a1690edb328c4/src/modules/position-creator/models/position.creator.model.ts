import { Field, ObjectType } from '@nestjs/graphql';
import { TransactionModel } from 'src/models/transaction.model';
import { SwapRouteModel } from 'src/modules/auto-router/models/auto-route.model';
import { DcdtTokenPaymentModel } from 'src/modules/tokens/models/dcdt.token.payment.model';

@ObjectType()
export class PositionCreatorModel {
    @Field()
    address: string;

    constructor(init: Partial<PositionCreatorModel>) {
        Object.assign(this, init);
    }
}

@ObjectType()
export class LiquidityPositionSingleTokenModel {
    @Field(() => DcdtTokenPaymentModel)
    payment: DcdtTokenPaymentModel;

    @Field(() => [SwapRouteModel])
    swaps: SwapRouteModel[];

    @Field(() => [TransactionModel], { nullable: true })
    transactions?: TransactionModel[];

    constructor(init: Partial<LiquidityPositionSingleTokenModel>) {
        Object.assign(this, init);
    }
}

@ObjectType()
export class FarmPositionSingleTokenModel {
    @Field(() => DcdtTokenPaymentModel)
    payment: DcdtTokenPaymentModel;

    @Field(() => [SwapRouteModel])
    swaps: SwapRouteModel[];

    @Field(() => [TransactionModel], { nullable: true })
    transactions?: TransactionModel[];

    constructor(init: Partial<FarmPositionSingleTokenModel>) {
        Object.assign(this, init);
    }
}

@ObjectType()
export class DualFarmPositionSingleTokenModel {
    @Field(() => DcdtTokenPaymentModel)
    payment: DcdtTokenPaymentModel;

    @Field(() => [SwapRouteModel])
    swaps: SwapRouteModel[];

    @Field(() => [TransactionModel], { nullable: true })
    transactions?: TransactionModel[];

    constructor(init: Partial<DualFarmPositionSingleTokenModel>) {
        Object.assign(this, init);
    }
}

@ObjectType()
export class StakingPositionSingleTokenModel {
    @Field(() => DcdtTokenPaymentModel)
    payment: DcdtTokenPaymentModel;

    @Field(() => [SwapRouteModel])
    swaps: SwapRouteModel[];

    @Field(() => [TransactionModel], { nullable: true })
    transactions?: TransactionModel[];

    constructor(init: Partial<StakingPositionSingleTokenModel>) {
        Object.assign(this, init);
    }
}

@ObjectType()
export class EnergyPositionSingleTokenModel {
    @Field(() => DcdtTokenPaymentModel)
    payment: DcdtTokenPaymentModel;

    @Field(() => [SwapRouteModel])
    swaps: SwapRouteModel[];

    @Field(() => [TransactionModel], { nullable: true })
    transactions?: TransactionModel[];

    constructor(init: Partial<EnergyPositionSingleTokenModel>) {
        Object.assign(this, init);
    }
}
