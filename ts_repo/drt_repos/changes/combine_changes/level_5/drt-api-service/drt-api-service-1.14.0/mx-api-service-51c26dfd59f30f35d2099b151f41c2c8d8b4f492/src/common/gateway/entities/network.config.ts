export class NetworkConfig {
  constructor(init?: Partial<NetworkConfig>) {
    Object.assign(this, init);
  }

  drt_adaptivity: boolean = false;
  drt_chain_id: string = '';
  drt_denomination: number = 0;
  drt_gas_per_data_byte: number = 0;
  drt_gas_price_modifier: string = '0';
  drt_hysteresis: string = '0';
  drt_latest_tag_software_version: string = '0';
  drt_max_gas_per_transaction: number = 0;
  drt_meta_consensus_group_size: number = 0;
  drt_min_gas_limit: number = 0;
  drt_min_gas_price: number = 0;
  drt_min_transaction_version: number = 0;
  drt_num_metachain_nodes: number = 0;
  drt_num_nodes_in_shard: number = 0;
  drt_num_shards_without_meta: number = 0;
  drt_rewards_top_up_gradient_point: string = '0';
  drt_round_duration: number = 0;
  drt_rounds_per_epoch: number = 0;
  drt_shard_consensus_group_size: number = 0;
  drt_start_time: number = 0;
  drt_top_up_factor: string = '0';
}
