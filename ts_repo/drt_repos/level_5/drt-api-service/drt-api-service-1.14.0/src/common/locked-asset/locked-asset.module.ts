import { Module } from '@nestjs/common';
import { LockedAssetService } from './locked-asset.service';
import { VmQueryModule } from '../../endpoints/vm.query/vm.query.module';
import { MoaModule } from 'src/endpoints/moa/moa.module';

@Module({
  imports: [
    VmQueryModule,
    MoaModule.forRoot(),
  ],
  providers: [
    LockedAssetService,
  ],
  exports: [
    LockedAssetService,
  ],
})
export class LockedAssetModule { }
