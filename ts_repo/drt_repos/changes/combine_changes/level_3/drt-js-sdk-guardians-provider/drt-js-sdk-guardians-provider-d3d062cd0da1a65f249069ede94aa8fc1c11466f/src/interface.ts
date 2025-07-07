export interface IAPIAccountResponse {
  address: string;
  balance: string;
  nonce: number;
  shard: number;
  rootHash: string;
  txCount: number;
  scrCount: number;
  developerReward: string;
  activeGuardianServiceUid?: string;
  activeGuardianActivationEpoch?: number;
  activeGuardianAddress?: string;
  pendingGuardianActivationEpoch?: number;
  pendingGuardianServiceUid?: string;
  pendingGuardianAddress?: string;
  isGuarded?: boolean;
}

export interface IInitData {
  activeGuardianServiceUid: string;
  isGuarded: boolean;
  activeGuardianAddress: string;
  pendingGuardianActivationEpoch: number;
  pendingGuardianAddress: string;
  address: string;
  apiAddress: string;
  networkId: string;
  registrationDelay: number;
  backoffWrongCode: number;
  pendingGuardianServiceUid?: string;
}

export interface IRegisterOptions {
  tag?: string;
}
