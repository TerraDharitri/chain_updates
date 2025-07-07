import { DcdtLockedAccount } from "./dcdt.locked.account";
import { ApiProperty } from "@nestjs/swagger";
import { SwaggerUtils } from "@terradharitri/sdk-nestjs-common";

export class DcdtSupply {
  constructor(init?: Partial<DcdtSupply>) {
    Object.assign(this, init);
  }

  @ApiProperty(SwaggerUtils.amountPropertyOptions())
  totalSupply: string = '0';

  @ApiProperty(SwaggerUtils.amountPropertyOptions())
  circulatingSupply: string = '0';

  @ApiProperty(SwaggerUtils.amountPropertyOptions())
  minted: string = '0';

  @ApiProperty(SwaggerUtils.amountPropertyOptions())
  burned: string = '0';

  @ApiProperty(SwaggerUtils.amountPropertyOptions())
  initialMinted: string = '0';

  @ApiProperty()
  lockedAccounts: DcdtLockedAccount[] | undefined = undefined;
}
