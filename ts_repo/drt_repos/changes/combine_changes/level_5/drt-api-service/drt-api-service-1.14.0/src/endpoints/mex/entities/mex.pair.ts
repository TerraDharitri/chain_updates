import { ApiProperty } from "@nestjs/swagger";
import { MoaPairState } from "./moa.pair.state";
import { MoaPairType } from "./moa.pair.type";
import { MoaPairExchange } from "./moa.pair.exchange";

export class MoaPair {
  constructor(init?: Partial<MoaPair>) {
    Object.assign(this, init);
  }

  @ApiProperty()
  address: string = '';

  @ApiProperty()
  id: string = '';

  @ApiProperty()
  symbol: string = '';

  @ApiProperty()
  name: string = '';

  @ApiProperty()
  price: number = 0;

  @ApiProperty()
  basePrevious24hPrice: number = 0;

  @ApiProperty()
  quotePrevious24hPrice: number = 0;

  @ApiProperty({ type: String, example: 'MOA-455c57' })
  baseId: string = '';

  @ApiProperty({ type: String, example: 'MOA' })
  baseSymbol: string = '';

  @ApiProperty({ type: String, example: 'MOA' })
  baseName: string = '';

  @ApiProperty({ type: Number, example: 0.00020596180499578328 })
  basePrice: number = 0;

  @ApiProperty({ type: String, example: 'WREWA-bd4d79' })
  quoteId: string = '';

  @ApiProperty({ type: String, example: 'WREWA' })
  quoteSymbol: string = '';

  @ApiProperty({ type: String, example: 'WrappedREWA' })
  quoteName: string = '';

  @ApiProperty({ type: Number, example: 145.26032 })
  quotePrice: number = 0;

  @ApiProperty({ type: Number, example: '347667206.84174806' })
  totalValue: number = 0;

  @ApiProperty({ type: Number, example: '2109423.4531209776' })
  volume24h: number | undefined;

  @ApiProperty({ enum: MoaPairState })
  state: MoaPairState = MoaPairState.inactive;

  @ApiProperty({ enum: MoaPairType })
  type: MoaPairType = MoaPairType.experimental;

  @ApiProperty({ type: String, example: 'jungledex' })
  exchange: MoaPairExchange | undefined;

  @ApiProperty({ type: Number, nullable: true })
  tradesCount: number | undefined = undefined;

  @ApiProperty({ type: Number, nullable: true })
  tradesCount24h: number | undefined = undefined;

  @ApiProperty({ type: Number, nullable: true })
  deployedAt: number | undefined = undefined;

  @ApiProperty({ type: Boolean, nullable: true })
  hasFarms?: boolean;

  @ApiProperty({ type: Boolean, nullable: true })
  hasDualFarms?: boolean;
}
