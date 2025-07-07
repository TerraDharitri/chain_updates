import { ApiProperty } from "@nestjs/swagger";
export class MoaToken {
  constructor(init?: Partial<MoaToken>) {
    Object.assign(this, init);
  }

  @ApiProperty({ type: String, example: 'MOA-455c57' })
  id: string = '';

  @ApiProperty({ type: String, example: 'MOA' })
  symbol: string = '';

  @ApiProperty({ type: String, example: 'MOA' })
  name: string = '';

  @ApiProperty({ type: Number, example: 0.000206738758250580 })
  price: number = 0;

  @ApiProperty({ type: Number, example: 0.000206738758250580 })
  previous24hPrice: number = 0;

  @ApiProperty({ type: Number, example: 0.000206738758250580 })
  previous24hVolume: number | undefined = 0;

  @ApiProperty({ type: Number, nullable: true })
  tradesCount: number | undefined = 0;
}
