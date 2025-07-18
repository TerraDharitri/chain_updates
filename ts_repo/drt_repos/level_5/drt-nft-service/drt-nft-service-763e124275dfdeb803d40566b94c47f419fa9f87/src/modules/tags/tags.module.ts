import { Module } from '@nestjs/common';
import { TagsService } from './tags.service';
import { DrtCommunicationModule } from 'src/common';
import { TagsResolver } from './tags.resolver';

@Module({
  providers: [TagsService, TagsResolver],
  imports: [DrtCommunicationModule],
  exports: [TagsService],
})
export class TagsModuleGraph {}
