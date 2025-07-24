import { forwardRef, Module } from "@nestjs/common";
import { AssetsModule } from "src/common/assets/assets.module";
import { PluginModule } from "src/plugins/plugin.module";
import { NftMediaModule } from "src/queue.worker/nft.worker/queue/job-services/media/nft.media.module";
import { NftMetadataModule } from "src/queue.worker/nft.worker/queue/job-services/metadata/nft.metadata.module";
import { CollectionModule } from "../collections/collection.module";
import { DcdtModule } from "../dcdt/dcdt.module";
import { MoaModule } from "../moa/moa.module";
import { TokenModule } from "../tokens/token.module";
import { NftExtendedAttributesService } from "./nft.extendedattributes.service";
import { NftService } from "./nft.service";
import { LockedAssetModule } from "../../common/locked-asset/locked-asset.module";

@Module({
  imports: [
    forwardRef(() => DcdtModule),
    forwardRef(() => TokenModule),
    forwardRef(() => CollectionModule),
    forwardRef(() => PluginModule),
    forwardRef(() => NftMetadataModule),
    forwardRef(() => MoaModule.forRoot()),
    forwardRef(() => AssetsModule),
    forwardRef(() => LockedAssetModule),
    NftMediaModule,
  ],
  providers: [
    NftService, NftExtendedAttributesService,
  ],
  exports: [
    NftService, NftExtendedAttributesService,
  ],
})
export class NftModule { }
