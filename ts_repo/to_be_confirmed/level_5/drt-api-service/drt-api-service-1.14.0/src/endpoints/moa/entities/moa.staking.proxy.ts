import { ApiProperty } from "@nestjs/swagger";

export class MoaStakingProxy {
  constructor(init?: Partial<MoaStakingProxy>) {
    Object.assign(this, init);
  }

  @ApiProperty({ type: String, example: 'drt1qqqqqqqqqqqqqpgq2ymdr66nk5hx32j2tqdeqv9ajm9dj9uu2jpsv3u0ej' })
  address: string = '';

  @ApiProperty()
  dualYieldTokenName: string = '';

  @ApiProperty()
  dualYieldTokenCollection: string = '';

  static fromQueryResponse(response: any): MoaStakingProxy {
    const stakingProxy = new MoaStakingProxy();
    stakingProxy.address = response.address;
    stakingProxy.dualYieldTokenName = response.dualYieldToken.name;
    stakingProxy.dualYieldTokenCollection = response.dualYieldToken.collection;

    return stakingProxy;
  }
}
