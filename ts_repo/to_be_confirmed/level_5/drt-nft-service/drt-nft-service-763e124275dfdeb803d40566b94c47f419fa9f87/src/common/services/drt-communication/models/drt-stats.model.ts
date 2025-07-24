export class DrtStats {
  shards: number;
  blocks: number;
  accounts: number;
  transactions: number;
  refreshRate: number;
  epoch: number;
  roundsPassed: number;
  roundsPerEpoch: number;

  constructor(init?: Partial<DrtStats>) {
    Object.assign(this, init);
  }
}
