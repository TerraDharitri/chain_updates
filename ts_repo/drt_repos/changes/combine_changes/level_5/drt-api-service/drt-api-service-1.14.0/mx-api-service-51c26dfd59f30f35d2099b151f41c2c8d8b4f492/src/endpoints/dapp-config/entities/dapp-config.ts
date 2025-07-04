import { ApiProperty } from "@nestjs/swagger";

export class DappConfig {
  constructor(init?: Partial<DappConfig>) {
    Object.assign(this, init);
  }

  @ApiProperty({ type: String, example: 'mainnet' })
  id: string = '';

  @ApiProperty({ type: String, example: 'Mainnet' })
  name: string = '';

  @ApiProperty({ type: String, example: 'REWA' })
  rewaLabel: string = '';

  @ApiProperty({ type: String, example: '4' })
  decimals: string = '';

  @ApiProperty({ type: String, example: '18' })
  rewaDenomination: string = '';

  @ApiProperty({ type: String, example: '1500' })
  gasPerDataByte: string = '';

  @ApiProperty({ type: String, example: '4000' })
  apiTimeout: string = '';

  @ApiProperty({ type: String, example: 'https://durian.page.link/?apn=com.dharitri.durian.wallet&isi=1519405832&ibi=com.dharitri.durian.wallet&link=https://durian.com/' })
  walletConnectDeepLink: string = '';

  @ApiProperty({ type: [String], example: 'https://bridge.walletconnect.org' })
  walletConnectBridgeAddresses: string = '';

  @ApiProperty({ type: String, example: 'https://wallet.dharitri.org' })
  walletAddress: string = '';

  @ApiProperty({ type: String, example: 'https://api.dharitri.org' })
  apiAddress: string = '';

  @ApiProperty({ type: String, example: 'https://explorer.dharitri.org' })
  explorerAddress: string = '';

  @ApiProperty({ type: String, example: '1' })
  chainId: string = '';

  @ApiProperty({ type: Number, example: 6000 })
  refreshRate: number = 0;
}
