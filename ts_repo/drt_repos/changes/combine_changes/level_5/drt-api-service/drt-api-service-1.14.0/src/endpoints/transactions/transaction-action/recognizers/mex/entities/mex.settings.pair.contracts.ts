export class MoaSettingsPairContracts {
  constructor(init?: Partial<MoaSettingsPairContracts>) {
    Object.assign(this, init);
  }

  rewaMoa: string = '';
  rewaUsd: string = '';
}
