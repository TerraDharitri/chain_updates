import { ObjectType, Field, Int } from '@nestjs/graphql';
import { nestedFieldComplexity } from 'src/helpers/complexity/field.estimators';
import { DcdtToken } from 'src/modules/tokens/models/dcdtToken.model';

@ObjectType()
export class WrapModel {
    @Field()
    address: string;
    @Field(() => Int)
    shard: number;
    @Field({ complexity: nestedFieldComplexity })
    wrappedToken: DcdtToken;

    constructor(init?: Partial<WrapModel>) {
        Object.assign(this, init);
    }
}
