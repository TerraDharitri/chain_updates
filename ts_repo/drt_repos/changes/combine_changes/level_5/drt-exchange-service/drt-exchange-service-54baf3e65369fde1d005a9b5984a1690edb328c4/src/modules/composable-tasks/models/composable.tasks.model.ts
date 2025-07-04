import { EnumType, EnumVariantDefinition } from '@terradharitri/sdk-core/out';
import { Field, ObjectType } from '@nestjs/graphql';

export enum ComposableTaskType {
    WRAP_REWA = 'WrapREWA',
    UNWRAP_REWA = 'UnwrapREWA',
    SWAP = 'Swap',
    ROUTER_SWAP = 'RouterSwap',
    SEND_REWA_OR_DCDT = 'SendRewaOrDcdt',
}

export class ComposableTaskEnumType {
    static getEnumType(): EnumType {
        return new EnumType('ComposableTaskType', [
            new EnumVariantDefinition('WrapREWA', 0),
            new EnumVariantDefinition('UnwrapREWA', 1),
            new EnumVariantDefinition('Swap', 2),
            new EnumVariantDefinition('RouterSwap', 3),
            new EnumVariantDefinition('SendRewaOrDcdt', 4),
        ]);
    }
}

@ObjectType()
export class ComposableTaskModel {
    @Field()
    address: string;

    constructor(init: Partial<ComposableTaskModel>) {
        Object.assign(this, init);
    }
}
