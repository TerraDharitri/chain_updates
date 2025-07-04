import { ApiProperty } from "@nestjs/swagger";

export class MoaTokenChart {
  constructor(init?: Partial<MoaTokenChart>) {
    Object.assign(this, init);
  }

  @ApiProperty()
  timestamp: number = 0;

  @ApiProperty()
  value: number = 0;
}
