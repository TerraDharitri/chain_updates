import { ResolveField, Resolver } from '@nestjs/graphql';
import { DcdtToken, DcdtTokenType } from '../tokens/models/dcdtToken.model';
import { TokensResolver } from '../tokens/token.resolver';
import { UserToken } from './models/user.model';

@Resolver(() => UserToken)
export class UserTokenResolver extends TokensResolver {
    @ResolveField(() => String)
    async type(parent: DcdtToken): Promise<string> {
        if (parent.type !== DcdtTokenType.FungibleLpToken) {
            return super.type(parent);
        }
        return parent.type;
    }
}
