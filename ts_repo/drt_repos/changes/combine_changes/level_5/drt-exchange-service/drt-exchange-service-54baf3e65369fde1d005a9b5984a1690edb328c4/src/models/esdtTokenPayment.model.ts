import {
    BigUIntType,
    EnumType,
    EnumVariantDefinition,
    FieldDefinition,
    StructType,
    TokenIdentifierType,
    U64Type,
} from '@terradharitri/sdk-core';
import { ObjectType, Field, Int, InputType } from '@nestjs/graphql';

@ObjectType('DcdtTokenPayment')
@InputType('DcdtTokenPaymentInput')
export class DcdtTokenPayment {
    @Field(() => Int, { nullable: true })
    tokenType: number;
    @Field()
    tokenID: string;
    @Field(() => Int)
    nonce: number;
    @Field()
    amount: string;

    constructor(init?: Partial<DcdtTokenPayment>) {
        Object.assign(this, init);
    }
}

export class DcdtTokenPaymentStruct {
    tokenType: DcdtTokenType;
    tokenID: string;
    nonce: number;
    amount: string;

    constructor(init?: Partial<DcdtTokenPaymentStruct>) {
        Object.assign(this, init);
    }

    static getStructure(): StructType {
        return new StructType('DcdtTokenPayment', [
            new FieldDefinition('token_type', '', DcdtTokenType.getEnum()),
            new FieldDefinition(
                'token_identifier',
                '',
                new TokenIdentifierType(),
            ),
            new FieldDefinition('token_nonce', '', new U64Type()),
            new FieldDefinition('amount', '', new BigUIntType()),
        ]);
    }
}

export class DcdtTokenType {
    name: string;
    discriminant: number;

    static getEnum(): EnumType {
        return new EnumType(DcdtTokenType.name, [
            new EnumVariantDefinition('Fungible', 0),
            new EnumVariantDefinition('NonFungible', 1),
            new EnumVariantDefinition('SemiFungible', 2),
            new EnumVariantDefinition('Meta', 3),
            new EnumVariantDefinition('Invalid', 4),
        ]);
    }
}

export class RewaOrDcdtTokenPayment {
    tokenIdentifier: string;
    nonce: number;
    amount: string;

    constructor(init?: Partial<RewaOrDcdtTokenPayment>) {
        Object.assign(this, init);
    }

    static getStructure(): StructType {
        return new StructType('RewaOrDcdtTokenPayment', [
            new FieldDefinition(
                'token_identifier',
                '',
                new TokenIdentifierType(),
            ),
            new FieldDefinition('token_nonce', '', new U64Type()),
            new FieldDefinition('amount', '', new BigUIntType()),
        ]);
    }
}
