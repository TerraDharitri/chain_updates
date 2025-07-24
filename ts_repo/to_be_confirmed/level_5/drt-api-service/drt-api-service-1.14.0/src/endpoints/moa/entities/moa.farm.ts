import { ApiProperty } from "@nestjs/swagger";
import { MoaFarmType } from "./moa.farm.type";
import { MoaToken } from "./moa.token";

export class MoaFarm {
  constructor(init?: Partial<MoaFarm>) {
    Object.assign(this, init);
  }

  @ApiProperty({ enum: MoaFarmType })
  type: MoaFarmType = MoaFarmType.standard;

  @ApiProperty({ nullable: true, required: false })
  version?: string;

  @ApiProperty({ type: String, example: 'drt1qqqqqqqqqqqqqpgqzps75vsk97w9nsx2cenv2r2tyxl4fl402jpsmzscxv' })
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
  farmingId: string = '';

  @ApiProperty()
  farmingSymbol: string = '';

  @ApiProperty()
  farmingName: string = '';

  @ApiProperty()
  farmingPrice: number = 0;

  @ApiProperty()
  farmedId: string = '';

  @ApiProperty()
  farmedSymbol: string = '';

  @ApiProperty()
  farmedName: string = '';

  @ApiProperty()
  farmedPrice: number = 0;

  static fromFarmQueryResponse(response: any): MoaFarm {
    let price = Number(response.farmTokenPriceUSD);

    const symbol = response.farmToken.collection.split('-')[0];
    if (['REWAUSDCF', 'REWAUSDCFL'].includes(symbol)) {
      price = price / (10 ** 12) * 2;
    }

    const moaFarm = new MoaFarm();
    moaFarm.type = MoaFarmType.standard;
    moaFarm.version = response.version;
    moaFarm.address = response.address;
    moaFarm.id = response.farmToken.collection;
    moaFarm.symbol = symbol;
    moaFarm.name = response.farmToken.name;
    moaFarm.price = price;
    moaFarm.farmingId = response.farmingToken.identifier;
    moaFarm.farmingSymbol = response.farmingToken.identifier.split('-')[0];
    moaFarm.farmingName = response.farmingToken.name;
    moaFarm.farmingPrice = Number(response.farmingTokenPriceUSD);
    moaFarm.farmedId = response.farmedToken.identifier;
    moaFarm.farmedSymbol = response.farmedToken.identifier.split('-')[0];
    moaFarm.farmedName = response.farmedToken.name;
    moaFarm.farmedPrice = Number(response.farmedTokenPriceUSD);

    return moaFarm;
  }

  static fromStakingFarmResponse(response: any, pairs: Record<string, MoaToken>): MoaFarm {
    const price = pairs[response.farmingToken.identifier]?.price ?? 0;

    const moaFarm = new MoaFarm();
    moaFarm.type = MoaFarmType.metastaking;
    moaFarm.address = response.address;
    moaFarm.id = response.farmToken.collection;
    moaFarm.symbol = response.farmToken.collection.split('-')[0];
    moaFarm.name = response.farmToken.name;
    moaFarm.price = price;
    moaFarm.farmingId = response.farmingToken.identifier;
    moaFarm.farmingSymbol = response.farmingToken.identifier.split('-')[0];
    moaFarm.farmingName = response.farmingToken.name;
    moaFarm.farmingPrice = price;
    moaFarm.farmedId = response.farmingToken.identifier;
    moaFarm.farmedSymbol = response.farmingToken.identifier.split('-')[0];
    moaFarm.farmedName = response.farmingToken.name;
    moaFarm.farmedPrice = price;

    return moaFarm;
  }
}
