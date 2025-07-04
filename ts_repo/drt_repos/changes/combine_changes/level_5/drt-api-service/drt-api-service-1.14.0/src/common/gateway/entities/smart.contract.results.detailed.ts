import { SmartContractResult } from "src/endpoints/sc-results/entities/smart.contract.result";

export class GatewaySmartContractResults extends SmartContractResult {
  constructor(init?: Partial<GatewaySmartContractResults>) {
    super();
    Object.assign(this, init);
  }
  tokens: string[] = [];
  dcdtValues: string[] = [];
  receivers: string[] = [];
  receiversShardIDs: number[] = [];
  operation: string = '';
}
