import { Field, ObjectType, registerEnumType } from '@nestjs/graphql';
import { DcdtToken } from 'src/modules/tokens/models/dcdtToken.model';

export enum TradingActivityAction {
    'BUY' = 'BUY',
    'SELL' = 'SELL',
}

registerEnumType(TradingActivityAction, { name: 'TradingActivityAction' });

@ObjectType()
export class TradingActivityModel {
    @Field()
    hash: string;
    @Field()
    timestamp: string;
    @Field()
    action: TradingActivityAction;
    @Field()
    inputToken: DcdtToken;
    @Field()
    outputToken: DcdtToken;

    constructor(init?: Partial<TradingActivityModel>) {
        Object.assign(this, init);
    }
}
