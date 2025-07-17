import { ApiProperty } from "@nestjs/swagger";
import { DcdtType } from "../../dcdt/entities/dcdt.type";
import { DcdtSubType } from "src/endpoints/dcdt/entities/dcdt.sub.type";
import { TokenOwnersHistory } from "./token.owner.history";

export class TokenProperties {
  constructor(init?: Partial<TokenProperties>) {
    Object.assign(this, init);
  }

  @ApiProperty()
  identifier: string = '';

  @ApiProperty()
  name: string = '';

  @ApiProperty()
  type: DcdtType = DcdtType.NonFungibleDCDT;

  @ApiProperty()
  subType: DcdtSubType | undefined = undefined;

  @ApiProperty()
  owner: string = '';

  @ApiProperty()
  wiped: string = '';

  @ApiProperty()
  decimals: number = 0;

  @ApiProperty()
  isPaused: boolean = false;

  @ApiProperty()
  tags: string[] = [];

  @ApiProperty()
  royalties: number = 0;

  @ApiProperty()
  uris: string[] = [];

  @ApiProperty()
  url: string = '';

  @ApiProperty()
  canUpgrade: boolean = false;

  @ApiProperty()
  canMint: boolean = false;

  @ApiProperty()
  canBurn: boolean = false;

  @ApiProperty()
  canChangeOwner: boolean = false;

  @ApiProperty()
  canPause: boolean = false;

  @ApiProperty()
  canFreeze: boolean = false;

  @ApiProperty()
  canWipe: boolean = false;

  @ApiProperty()
  canAddSpecialRoles: boolean = false;

  @ApiProperty()
  canTransferNFTCreateRole: boolean = false;

  @ApiProperty()
  NFTCreateStopped: boolean = false;

  @ApiProperty()
  timestamp: number = 0;

  @ApiProperty()
  ownersHistory: TokenOwnersHistory[] = [];
}
