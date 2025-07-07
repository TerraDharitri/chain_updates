import { SwaggerUtils } from "@terradharitri/sdk-nestjs-common";
import { ApiProperty } from "@nestjs/swagger";

export class AccountUndelegation {
  constructor(init?: Partial<AccountUndelegation>) {
    Object.assign(this, init);
  }

  @ApiProperty(SwaggerUtils.amountPropertyOptions())
  amount: string = '';

  @ApiProperty({ type: Number })
  seconds: number = 0;
}
