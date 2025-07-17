export class NetworkEconomics {
  constructor(init?: Partial<NetworkEconomics>) {
    Object.assign(this, init);
  }

  drt_dev_rewards: string = '';
  drt_epoch_for_economics_data: number = 0;
  drt_inflation: string = '0';
  drt_total_base_staked_value: string = '0';
  drt_total_fees: string = '0';
  drt_total_supply: string = '0';
  drt_total_top_up_value: string = '0';
}
