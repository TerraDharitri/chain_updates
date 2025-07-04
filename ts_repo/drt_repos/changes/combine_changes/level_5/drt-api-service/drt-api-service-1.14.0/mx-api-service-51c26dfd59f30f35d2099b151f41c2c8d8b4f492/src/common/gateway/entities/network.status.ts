export class NetworkStatus {
  constructor(init?: Partial<NetworkStatus>) {
    Object.assign(this, init);
  }

  drt_cross_check_block_height: string = '';
  drt_current_round: number = 0;
  drt_epoch_number: number = 0;
  drt_highest_final_nonce: number = 0;
  drt_nonce: number = 0;
  drt_nonce_at_epoch_start: number = 0;
  drt_nonces_passed_in_current_epoch: number = 0;
  drt_round_at_epoch_start: number = 0;
  drt_rounds_passed_in_current_epoch: number = 0;
  drt_rounds_per_epoch: number = 0;
}
