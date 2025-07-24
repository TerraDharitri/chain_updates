import { ApiProperty } from '@nestjs/swagger';
import { MoaSettings } from './moa.settings';

export class MoaEconomics {
  constructor(init?: Partial<MoaEconomics>) {
    Object.assign(this, init);
  }

  @ApiProperty({ type: Number, example: 8045920000000 })
  totalSupply: number = 0;

  @ApiProperty({ type: Number, example: 4913924072690 })
  circulatingSupply: number = 0;

  @ApiProperty({ type: Number, example: 0.00020552146843751037 })
  price: number = 0;

  @ApiProperty({ type: Number, example: 1009916891 })
  marketCap: number = 0;

  @ApiProperty({ type: Number, example: 13680479 })
  volume24h: number = 0;

  @ApiProperty({ type: Number, example: 15 })
  marketPairs: number = 0;

  static fromQueryResponse(response: any, settings: MoaSettings): MoaEconomics {
    const totalSupply = 8_045_920_000_000;
    const price = Number(response.moaPriceUSD);
    const circulatingSupply = Number(response.moaSupply);
    const marketCap = Math.round(circulatingSupply * price);
    const volume24h = Math.round(Number(response.factory.totalVolumeUSD24h));
    const marketPairs = settings.pairContracts.length;

    return new MoaEconomics({
      totalSupply,
      price,
      circulatingSupply,
      marketCap,
      volume24h,
      marketPairs,
    });
  }
}
