import { ApiProperty } from '@nestjs/swagger';
import { MoaPairType } from './moa.pair.type';

export class MoaTokenType {
  constructor(init?: Partial<MoaTokenType>) {
    Object.assign(this, init);
  }

  @ApiProperty({ type: String, example: '' })
  identifier: string = '';

  @ApiProperty({ enum: MoaPairType })
  type: MoaPairType = MoaPairType.experimental;
}
