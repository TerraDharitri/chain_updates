import { SwaggerUtils } from "@terradharitri/sdk-nestjs-common";
import { ApiProperty } from "@nestjs/swagger";
import { Token } from "./token";
import { MoaPairType } from "src/endpoints/moa/entities/moa.pair.type";

export class TokenWithBalance extends Token {
  constructor(init?: Partial<TokenWithBalance>) {
    super();
    Object.assign(this, init);
  }

  @ApiProperty(SwaggerUtils.amountPropertyOptions())
  balance: string = '';

  @ApiProperty({ type: Number, nullable: true, required: false })
  valueUsd: number | undefined = undefined;

  @ApiProperty({ type: String, nullable: true, required: false })
  attributes: string | undefined = undefined;

  @ApiProperty({ enum: MoaPairType })
  drtPairType: MoaPairType = MoaPairType.experimental;
}
