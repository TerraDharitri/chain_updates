import { Parent, ResolveField, Resolver } from '@nestjs/graphql';
import { AssetsModel } from './models/assets.model';
import { DcdtToken } from './models/dcdtToken.model';

@Resolver(() => DcdtToken)
export class TokensResolver {

    @ResolveField(() => AssetsModel, { nullable: true })
    async assets(@Parent() parent: DcdtToken): Promise<AssetsModel> {
        return new AssetsModel(parent.assets);
    }
}
