export interface ScResult {
  scHash: string
  nonce: number;
  gasLimit: string;
  gasPrice: string;
  value: string;
  sender: string;
  receiver: string;
  senderShard: number;
  receiverShard: number;
  data: string;
  prevTxHash: string;
  originalTxHash: string;
  callType: string;
  timestamp: number;
  tokens: string[];
  dcdtValues: string[];
  operation: string;
}
