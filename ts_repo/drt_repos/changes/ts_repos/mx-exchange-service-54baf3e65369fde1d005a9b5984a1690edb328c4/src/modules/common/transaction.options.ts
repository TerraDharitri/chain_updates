import { TokenTransfer } from '@terradharitri/sdk-core';

export class TransactionOptions {
    chainID?: string;
    function: string;
    gasLimit: number;
    sender: string;
    arguments?: any[];
    tokenTransfers?: TokenTransfer[];
    nativeTransferAmount?: string;

    constructor(init?: Partial<TransactionOptions>) {
        Object.assign(this, init);
    }
}
