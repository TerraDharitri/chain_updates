import { ApiProperty } from "@nestjs/swagger";

export class NftOwner {
  constructor(init?: Partial<NftOwner>) {
    Object.assign(this, init);
  }

  @ApiProperty({ type: String, example: 'drt1qga7ze0l03chfgru0a32wxqf2226nzrxnyhzer9lmudqhjgy7ycq0wn4su' })
  address: string = '';

  @ApiProperty({ type: String, default: '1' })
  balance: string = '';
}
