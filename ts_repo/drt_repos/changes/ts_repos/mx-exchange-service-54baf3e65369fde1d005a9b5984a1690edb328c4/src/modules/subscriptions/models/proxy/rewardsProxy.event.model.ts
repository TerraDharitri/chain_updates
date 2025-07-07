import { RewardsProxyEvent } from '@terradharitri/sdk-exchange';
import { Address } from '@terradharitri/sdk-core';
import { Field, ObjectType } from '@nestjs/graphql';
import { GenericToken } from 'src/models/genericToken.model';
import { WrappedFarmTokenAttributesModel } from 'src/modules/proxy/models/wrappedFarmTokenAttributes.model';
import { GenericEventModel } from '../generic.event.model';

@ObjectType()
export class RewardsProxyEventModel extends GenericEventModel {
    @Field(() => String)
    private farmAddress: Address;
    @Field(() => GenericToken)
    private oldWrappedFarmToken: GenericToken;
    @Field(() => GenericToken)
    private newWrappedFarmToken: GenericToken;
    @Field(() => WrappedFarmTokenAttributesModel)
    private oldWrappedFarmAttributes: WrappedFarmTokenAttributesModel;
    @Field(() => WrappedFarmTokenAttributesModel)
    private newWrappedFarmAttributes: WrappedFarmTokenAttributesModel;
    @Field()
    private createdWithMerge: boolean;

    constructor(init?: Partial<RewardsProxyEvent>) {
        super(init);
        Object.assign(this, init);
    }
}
