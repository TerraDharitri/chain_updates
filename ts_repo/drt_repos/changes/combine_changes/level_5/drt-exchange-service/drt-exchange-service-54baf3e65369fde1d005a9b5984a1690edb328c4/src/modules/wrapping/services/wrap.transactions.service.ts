import { Injectable } from '@nestjs/common';
import { Token, TokenTransfer } from '@terradharitri/sdk-core';
import { TransactionModel } from '../../../models/transaction.model';
import { drtConfig, gasConfig } from '../../../config';
import { WrapService } from './wrap.service';
import { MXProxyService } from '../../../services/dharitri-communication/drt.proxy.service';
import { TransactionOptions } from 'src/modules/common/transaction.options';

@Injectable()
export class WrapTransactionsService {
    constructor(
        private readonly drtProxy: MXProxyService,
        private readonly wrapService: WrapService,
    ) {}

    async wrapRewa(sender: string, amount: string): Promise<TransactionModel> {
        const shardID = await this.drtProxy.getAddressShardID(sender);

        return this.drtProxy.getWrapSmartContractTransaction(
            shardID,
            new TransactionOptions({
                function: 'wrapRewa',
                chainID: drtConfig.chainID,
                gasLimit: gasConfig.wrapREWA,
                sender: sender,
                nativeTransferAmount: amount,
            }),
        );
    }

    async unwrapRewa(
        sender: string,
        amount: string,
    ): Promise<TransactionModel> {
        const [shardID, wrappedRewaToken] = await Promise.all([
            this.drtProxy.getAddressShardID(sender),
            this.wrapService.wrappedRewaToken(),
        ]);

        return this.drtProxy.getWrapSmartContractTransaction(
            shardID,
            new TransactionOptions({
                function: 'unwrapRewa',
                chainID: drtConfig.chainID,
                gasLimit: gasConfig.wrapREWA,
                sender: sender,
                tokenTransfers: [
                    new TokenTransfer({
                        token: new Token({
                            identifier: wrappedRewaToken.identifier,
                        }),
                        amount: BigInt(amount),
                    }),
                ],
            }),
        );
    }
}
