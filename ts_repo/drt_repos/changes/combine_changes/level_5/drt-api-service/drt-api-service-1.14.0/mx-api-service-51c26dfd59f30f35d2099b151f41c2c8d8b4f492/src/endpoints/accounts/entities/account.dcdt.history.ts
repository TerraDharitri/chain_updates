import { ApiProperty } from "@nestjs/swagger";
import { AccountHistory } from "./account.history";

export class AccountDcdtHistory extends AccountHistory {
  constructor(init?: Partial<AccountDcdtHistory>) {
    super();
    Object.assign(this, init);
  }

  @ApiProperty({ type: String, example: 'WREWA-bd4d79' })
  token: string = '';

  @ApiProperty({ type: String, example: 'XPACHIEVE-5a0519-01' })
  identifier: string | undefined = undefined;
}
