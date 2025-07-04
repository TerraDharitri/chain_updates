import { ObjectType } from '@nestjs/graphql';
import { AssetsModel } from './assets.model';
import { IDcdtToken } from './dcdtToken.interface';

@ObjectType({
    implements: () => [IDcdtToken],
})
export class DcdtToken implements IDcdtToken {
    identifier: string;
    name: string;
    ticker: string;
    owner: string;
    minted?: string;
    burnt?: string;
    initialMinted?: string;
    decimals: number;
    price?: string;
    supply?: string;
    circulatingSupply?: string;
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
    type?: string;
    balance?: string;

    constructor(init?: Partial<DcdtToken>) {
        Object.assign(this, init);
        this.assets = new AssetsModel(init.assets);
    }
}
