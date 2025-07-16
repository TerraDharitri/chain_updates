import { Field, ObjectType } from '@nestjs/graphql';

@ObjectType()
export class DcdtTokenPaymentModel {
    @Field()
    tokenIdentifier: string;
    @Field()
    tokenNonce: number;
    @Field()
    amount: string;

    constructor(init: DcdtTokenPaymentModel) {
        Object.assign(this, init);
    }
}
