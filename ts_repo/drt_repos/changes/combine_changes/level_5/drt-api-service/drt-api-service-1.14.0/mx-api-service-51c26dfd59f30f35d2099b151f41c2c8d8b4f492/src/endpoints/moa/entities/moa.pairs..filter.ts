import { MoaPairExchange } from "./moa.pair.exchange";

export class MoaPairsFilter {
  constructor(init?: Partial<MoaPairsFilter>) {
    Object.assign(this, init);
  }
  exchange?: MoaPairExchange;
  includeFarms?: boolean;
}
