import { Field, ObjectType } from '@nestjs/graphql';
import { AssetsModel } from './assets.model';
import { IDcdtToken } from './dcdtToken.interface';
import { RolesModel } from './roles.model';
import { nestedFieldComplexity } from 'src/helpers/complexity/field.estimators';

export enum DcdtTokenType {
    FungibleToken = 'FungibleDCDT',
    FungibleLpToken = 'FungibleDCDT-LP',
}

export class BaseDcdtToken {
    identifier: string;
    decimals: number;

    constructor(init?: Partial<BaseDcdtToken>) {
        Object.assign(this, init);
    }

    static toDcdtToken(baseDcdtToken: BaseDcdtToken): DcdtToken {
        return new DcdtToken(baseDcdtToken);
    }
}

@ObjectType({
    implements: () => [IDcdtToken],
})
export class DcdtToken extends BaseDcdtToken implements IDcdtToken {
    name: string;
    ticker: string;
    owner: string;
    minted?: string;
    burnt?: string;
    initialMinted?: string;
    derivedREWA: string;
    price?: string;
    previous24hPrice?: string;
    previous7dPrice?: string;
    volumeUSD24h?: string;
    previous24hVolume?: string;
    liquidityUSD?: string;
    swapCount24h?: number;
    previous24hSwapCount?: number;
    trendingScore?: string;
    supply?: string;
    circulatingSupply?: string;
    @Field(() => AssetsModel, {
        nullable: true,
        complexity: nestedFieldComplexity,
    })
    assets?: AssetsModel;
    transactions: number;
    accounts: number;
    isPaused: boolean;
    canUpgrade: boolean;
    canMint: boolean;
    canBurn: boolean;
    canChangeOwner: boolean;
    canPause: boolean;
    canFreeze: boolean;
    canWipe: boolean;
    @Field(() => RolesModel, {
        nullable: true,
        complexity: nestedFieldComplexity,
    })
    roles?: RolesModel;
    type?: string;
    balance?: string;
    createdAt?: string;

    constructor(init?: Partial<DcdtToken>) {
        super(init);
        Object.assign(this, init);
        if (init.assets) {
            this.assets = new AssetsModel(init.assets);
        }
        if (init.roles) {
            this.roles = new RolesModel(init.roles);
        }
    }
}
