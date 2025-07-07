import {
    BinaryCodec,
    FieldDefinition,
    StructType,
    U64Type,
} from '@terradharitri/sdk-core';
import { DcdtTokenPayment } from '../dcdt-token-payment';
import { UnstakePairType } from './unstake.pair.type';

export class UnstakePair {
    readonly unlockEpoch: number;
    readonly lockedTokens: DcdtTokenPayment;
    readonly unlockedTokens: DcdtTokenPayment;

    constructor(init: UnstakePairType) {
        this.unlockEpoch = init.unlockEpoch;
        this.lockedTokens = new DcdtTokenPayment(init.lockedTokens);
        this.unlockedTokens = new DcdtTokenPayment(init.unlockedTokens);
    }

    toJSON(): UnstakePairType {
        return {
            unlockEpoch: this.unlockEpoch,
            lockedTokens: this.lockedTokens.toJSON(),
            unlockedTokens: this.unlockedTokens.toJSON(),
        };
    }

    static fromDecodedAttributes(decodedAttributes: any): UnstakePair {
        return new UnstakePair({
            unlockEpoch: decodedAttributes.unlock_epoch.toNumber(),
            lockedTokens: DcdtTokenPayment.fromDecodedAttributes(
                decodedAttributes.locked_tokens,
            ),
            unlockedTokens: DcdtTokenPayment.fromDecodedAttributes(
                decodedAttributes.unlocked_tokens,
            ),
        });
    }

    static fromAttributes(attributes: string): UnstakePair {
        const attributesBuffer = Buffer.from(attributes, 'base64');
        const codec = new BinaryCodec();

        const structType = this.getStructure();
        const [decoded] = codec.decodeNested(attributesBuffer, structType);

        return this.fromDecodedAttributes(decoded.valueOf());
    }

    static getStructure(): StructType {
        return new StructType('UnstakePair', [
            new FieldDefinition('unlock_epoch', '', new U64Type()),
            new FieldDefinition(
                'locked_tokens',
                '',
                DcdtTokenPayment.getStructure(),
            ),
            new FieldDefinition(
                'unlocked_tokens',
                '',
                DcdtTokenPayment.getStructure(),
            ),
        ]);
    }
}
