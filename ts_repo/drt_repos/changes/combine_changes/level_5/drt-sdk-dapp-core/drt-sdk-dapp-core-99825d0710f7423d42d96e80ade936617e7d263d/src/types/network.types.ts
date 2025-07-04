export interface NetworkType {
  id: string;
  chainId: string;
  name: string;
  rewaLabel: string;
  decimals: string;
  digits: string;
  gasPerDataByte: string;
  walletAddress: string;
  apiAddress: string;
  explorerAddress: string;
  apiTimeout: string;
  xAliasAddress?: string;
  roundDuration: number;
  iframeWalletAddress?: string;
  websocketUrl?: string;
  gasStationMetadata?: {
    fast: number;
    faster: number;
    excellentJustLikeMoveBalance?: number;
  }[];
}

export type CustomNetworkType = {
  [P in keyof NetworkType]?: NetworkType[P];
} & {
  /**
   * If set to `true`, network configuration in init app will prevent a call to `dapp/config`.
   */
  skipFetchFromServer?: boolean;
};

export interface ApiNetworkConfigType {
  drt_chain_id: string;
  drt_denomination: number;
  drt_gas_per_data_byte: number;
  drt_gas_price_modifier: string;
  drt_latest_tag_software_version: string;
  drt_max_gas_per_transaction: number;
  drt_meta_consensus_group_size: number;
  drt_min_gas_limit: number;
  drt_min_gas_price: number;
  drt_min_transaction_version: number;
  drt_num_metachain_nodes: number;
  drt_num_nodes_in_shard: number;
  drt_num_shards_without_meta: number;
  drt_rewards_top_up_gradient_point: string;
  drt_round_duration: number;
  drt_rounds_per_epoch: number;
  drt_shard_consensus_group_size: number;
  drt_start_time: number;
  drt_top_up_factor: string;
}
