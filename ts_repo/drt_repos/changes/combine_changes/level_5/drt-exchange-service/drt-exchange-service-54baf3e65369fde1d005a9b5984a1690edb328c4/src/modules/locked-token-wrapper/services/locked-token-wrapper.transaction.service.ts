import { Injectable } from '@nestjs/common';
import { MXProxyService } from '../../../services/dharitri-communication/drt.proxy.service';
import { TransactionModel } from '../../../models/transaction.model';
import { Token, TokenTransfer } from '@terradharitri/sdk-core';
import { gasConfig } from '../../../config';
import { InputTokenModel } from '../../../models/inputToken.model';
import { TransactionOptions } from 'src/modules/common/transaction.options';

@Injectable()
export class LockedTokenWrapperTransactionService {
    constructor(private readonly drtProxy: MXProxyService) {}

    async unwrapLockedToken(
        sender: string,
        inputToken: InputTokenModel,
    ): Promise<TransactionModel> {
        return this.drtProxy.getLockedTokenWrapperSmartContractTransaction(
            new TransactionOptions({
                sender: sender,
                gasLimit: gasConfig.lockedTokenWrapper.unwrapLockedToken,
                function: 'unwrapLockedToken',
                tokenTransfers: [
                    new TokenTransfer({
                        token: new Token({
                            identifier: inputToken.tokenID,
                            nonce: BigInt(inputToken.nonce),
                        }),
                        amount: BigInt(inputToken.amount),
                    }),
                ],
            }),
        );
    }

    async wrapLockedToken(
        sender: string,
        inputToken: InputTokenModel,
    ): Promise<TransactionModel> {
        return this.drtProxy.getLockedTokenWrapperSmartContractTransaction(
            new TransactionOptions({
                sender: sender,
                gasLimit: gasConfig.lockedTokenWrapper.wrapLockedToken,
                function: 'wrapLockedToken',
                tokenTransfers: [
                    new TokenTransfer({
                        token: new Token({
                            identifier: inputToken.tokenID,
                            nonce: BigInt(inputToken.nonce),
                        }),
                        amount: BigInt(inputToken.amount),
                    }),
                ],
            }),
        );
    }
}
