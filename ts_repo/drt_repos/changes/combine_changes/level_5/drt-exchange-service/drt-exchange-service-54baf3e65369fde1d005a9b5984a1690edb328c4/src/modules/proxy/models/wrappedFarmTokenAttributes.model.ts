import { Field, Int, ObjectType } from '@nestjs/graphql';
import { LockedAssetAttributesModel } from 'src/modules/locked-asset-factory/models/locked-asset.model';
import { DcdtTokenPaymentModel } from 'src/modules/tokens/models/dcdt.token.payment.model';
import { FarmTokenAttributesUnion } from '../../farm/models/farmTokenAttributes.model';
import { LockedAssetAttributesUnion } from './locked.assets.attributes.union';
import {
    WrappedLpTokenAttributesModel,
    WrappedLpTokenAttributesModelV2,
} from './wrappedLpTokenAttributes.model';

@ObjectType()
export class WrappedFarmTokenAttributesModel {
    @Field()
    identifier: string;
    @Field()
    attributes: string;
    @Field()
    farmTokenID: string;
    @Field(() => Int)
    farmTokenNonce: number;
    @Field()
    farmTokenAmount: string;
    @Field()
    farmTokenIdentifier: string;
    @Field(() => FarmTokenAttributesUnion)
    farmTokenAttributes: typeof FarmTokenAttributesUnion;
    @Field()
    farmingTokenID: string;
    @Field(() => Int)
    farmingTokenNonce: number;
    @Field()
    farmingTokenAmount: string;
    @Field(() => LockedAssetAttributesModel, { nullable: true })
    lockedAssetsAttributes: LockedAssetAttributesModel;
    @Field(() => WrappedLpTokenAttributesModel, { nullable: true })
    lockedLpProxyTokenAttributes: WrappedLpTokenAttributesModel;

    constructor(init?: Partial<WrappedFarmTokenAttributesModel>) {
        Object.assign(this, init);
    }
}

@ObjectType()
export class WrappedFarmTokenAttributesModelV2 {
    @Field()
    identifier: string;
    @Field()
    attributes: string;
    @Field(() => DcdtTokenPaymentModel)
    farmToken: DcdtTokenPaymentModel;
    @Field(() => DcdtTokenPaymentModel)
    proxyFarmingToken: DcdtTokenPaymentModel;
    @Field(() => FarmTokenAttributesUnion)
    farmTokenAttributes: typeof FarmTokenAttributesUnion;
    @Field(() => WrappedLpTokenAttributesModelV2, { nullable: true })
    lockedLpProxyTokenAttributes: WrappedLpTokenAttributesModelV2;

    constructor(init?: Partial<WrappedFarmTokenAttributesModelV2>) {
        Object.assign(this, init);
    }
}
