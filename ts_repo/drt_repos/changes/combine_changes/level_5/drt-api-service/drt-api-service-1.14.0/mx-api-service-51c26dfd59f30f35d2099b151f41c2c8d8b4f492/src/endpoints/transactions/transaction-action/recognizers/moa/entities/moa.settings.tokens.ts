export class MoaSettingsToken {
  constructor(init?: Partial<MoaSettingsToken>) {
    Object.assign(this, init);
  }

  wrewa: string = '';
  moa: string = '';
  busd: string = '';
  rewaMoa: string = '';
  rewaUsd: string = '';

  moaFarm: string = '';
  rewaMoaFarm: string = '';
  rewaUsdFarm: string = '';
}
