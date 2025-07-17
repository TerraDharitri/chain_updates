import { forwardRef, Module } from "@nestjs/common";
import { VmQueryModule } from "src/endpoints/vm.query/vm.query.module";
import { DcdtService } from "./dcdt.service";
import { TokenModule } from "../tokens/token.module";
import { DcdtAddressService } from "./dcdt.address.service";
import { NftModule } from "../nfts/nft.module";
import { CollectionModule } from "../collections/collection.module";
import { TransactionModule } from "../transactions/transaction.module";
import { MoaModule } from "../moa/moa.module";
import { AssetsModule } from "src/common/assets/assets.module";


@Module({
  imports: [
    forwardRef(() => NftModule),
    forwardRef(() => CollectionModule),
    forwardRef(() => TokenModule),
    VmQueryModule,
    forwardRef(() => TransactionModule),
    forwardRef(() => MoaModule.forRoot()),
    forwardRef(() => AssetsModule),
  ],
  providers: [
    DcdtService, DcdtAddressService,
  ],
  exports: [
    DcdtService, DcdtAddressService,
  ],
})
export class DcdtModule { }
