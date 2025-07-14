export class MoaSettingsFarmContracts {
  constructor(init?: Partial<MoaSettingsFarmContracts>) {
    Object.assign(this, init);
  }

  rewaMoa: string = '';
  rewaUsd: string = '';
  moa: string = '';
  exit: string = '';
}
