import { DcdtType } from "../../dcdt/entities/dcdt.type";

export class TokenTransferProperties {
  constructor(init?: Partial<TokenTransferProperties>) {
    Object.assign(this, init);
  }

  type: DcdtType = DcdtType.FungibleDCDT;
  token?: string;
  collection?: string;
  identifier?: string;
  ticker: string = '';
  decimals?: number = 0;
  name: string = '';
  svgUrl: string = '';
  valueUsd?: number;
  valueRewa?: number;
}
