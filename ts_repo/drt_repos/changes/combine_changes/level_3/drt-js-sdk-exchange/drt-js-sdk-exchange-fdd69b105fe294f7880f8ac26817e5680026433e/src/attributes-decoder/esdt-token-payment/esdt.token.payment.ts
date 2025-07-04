import {
    BigUIntType,
    BinaryCodec,
    EnumType,
    EnumVariantDefinition,
    FieldDefinition,
    StructType,
    TokenIdentifierType,
    U64Type,
} from '@terradharitri/sdk-core';
import { DcdtTokenPaymentType } from './dcdt.token.payment.type';

export class DcdtTokenPayment {
    tokenIdentifier: string;
    tokenNonce: number;
    amount: string;

    constructor(init: DcdtTokenPaymentType) {
        this.tokenIdentifier = init.tokenIdentifier;
        this.tokenNonce = init.tokenNonce;
        this.amount = init.amount;
    }

    toJSON(): DcdtTokenPaymentType {
        return {
            tokenIdentifier: this.tokenIdentifier,
            tokenNonce: this.tokenNonce,
            amount: this.amount,
        };
    }

    static fromAttributes(attributes: string): DcdtTokenPayment {
        const attributesBuffer = Buffer.from(attributes, 'base64');
        const codec = new BinaryCodec();

        const firstFourBytes = attributesBuffer.subarray(0, 4);

        const withTokenType = firstFourBytes[3] === 0 || firstFourBytes[0] > 0;

        const structType = this.getStructure(withTokenType);
        const [decoded] = codec.decodeNested(attributesBuffer, structType);
        return this.fromDecodedAttributes(decoded.valueOf());
    }

    static fromDecodedAttributes(decodedAttributes: any): DcdtTokenPayment {
        return new DcdtTokenPayment({
            tokenIdentifier: decodedAttributes.token_identifier.toString(),
            tokenNonce: decodedAttributes.token_nonce.toNumber(),
            amount: decodedAttributes.amount.toFixed(),
        });
    }

    static getStructure(withTokenType = false): StructType {
        const fieldDefinitions = [
            new FieldDefinition('token_type', '', DcdtTokenType.getEnum()),
            new FieldDefinition(
                'token_identifier',
                '',
                new TokenIdentifierType(),
            ),
            new FieldDefinition('token_nonce', '', new U64Type()),
            new FieldDefinition('amount', '', new BigUIntType()),
        ];

        return new StructType(
            'DcdtTokenPayment',
            withTokenType ? fieldDefinitions : fieldDefinitions.slice(1),
        );
    }
}

export class DcdtTokenType {
    name: string;
    discriminant: number;

    constructor(init: DcdtTokenType) {
        this.name = init.name;
        this.discriminant = init.discriminant;
    }

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
