import { BinaryCodec, FieldDefinition, StructType } from '@terradharitri/sdk-core';
import { DcdtTokenPayment } from '../dcdt-token-payment/dcdt.token.payment';
import { WrappedFarmTokenAttributesTypeV2 } from './proxy.token.types';

export class WrappedFarmTokenAttributesV2 {
    farmToken: DcdtTokenPayment;
    proxyFarmingToken: DcdtTokenPayment;

    constructor(init: WrappedFarmTokenAttributesTypeV2) {
        this.farmToken = new DcdtTokenPayment(init.farmToken);
        this.proxyFarmingToken = new DcdtTokenPayment(init.proxyFarmingToken);
    }

    toJSON(): WrappedFarmTokenAttributesTypeV2 {
        return {
            farmToken: this.farmToken.toJSON(),
            proxyFarmingToken: this.proxyFarmingToken.toJSON(),
        };
    }

    static fromAttributes(attributes: string): WrappedFarmTokenAttributesV2 {
        const attributesBuffer = Buffer.from(attributes, 'base64');
        const codec = new BinaryCodec();
        const structType = this.getStructure();
        const [decoded] = codec.decodeNested(attributesBuffer, structType);
        return this.fromDecodedAttributes(decoded.valueOf());
    }

    static fromDecodedAttributes(
        decodedAttributes: any,
    ): WrappedFarmTokenAttributesV2 {
        return new WrappedFarmTokenAttributesV2({
            farmToken: DcdtTokenPayment.fromDecodedAttributes(
                decodedAttributes.farm_token,
            ),
            proxyFarmingToken: DcdtTokenPayment.fromDecodedAttributes(
                decodedAttributes.proxy_farming_token,
            ),
        });
    }

    static getStructure(): StructType {
        return new StructType('WrappedLpTokenAttributes', [
            new FieldDefinition(
                'farm_token',
                '',
                DcdtTokenPayment.getStructure(),
            ),
            new FieldDefinition(
                'proxy_farming_token',
                '',
                DcdtTokenPayment.getStructure(),
            ),
        ]);
    }
}
