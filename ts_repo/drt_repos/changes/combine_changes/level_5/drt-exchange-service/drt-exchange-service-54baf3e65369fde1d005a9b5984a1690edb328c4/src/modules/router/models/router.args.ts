import { EnumType, EnumVariantDefinition } from '@terradharitri/sdk-core';
import { Field, InputType, registerEnumType } from '@nestjs/graphql';

@InputType()
export class SetLocalRoleOwnerArgs {
    @Field()
    tokenID: string;

    @Field()
    address: string;

    @Field(() => [DcdtLocalRole])
    roles: DcdtLocalRole[];
}

export enum DcdtLocalRole {
    None,
    Mint,
    Burn,
    NftCreate,
    NftAddQuantity,
    NftBurn,
}
registerEnumType(DcdtLocalRole, { name: 'DcdtLocalRole' });
export const DcdtLocalRoleEnumType = new EnumType('DcdtLocalRole', [
    new EnumVariantDefinition('None', 0),
    new EnumVariantDefinition('Mint', 1),
    new EnumVariantDefinition('Burn', 2),
    new EnumVariantDefinition('NftCreate', 3),
    new EnumVariantDefinition('NftAddQuantity', 4),
    new EnumVariantDefinition('NftBurn', 5),
]);
