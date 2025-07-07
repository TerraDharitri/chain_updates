import { DcdtTokenPayment } from '@terradharitri/sdk-exchange';
import { Injectable } from '@nestjs/common';
import { MXProxyService } from 'src/services/dharitri-communication/drt.proxy.service';
import {
    ComposableTaskEnumType,
    ComposableTaskType,
} from '../models/composable.tasks.model';
import { TransactionModel } from 'src/models/transaction.model';
import {
    BigUIntValue,
    BytesType,
    BytesValue,
    EnumValue,
    EnumVariantDefinition,
    Field,
    List,
    ListType,
    Struct,
    Token,
    TokenIdentifierValue,
    TokenTransfer,
    TypedValue,
    U64Value,
    VariadicValue,
} from '@terradharitri/sdk-core';
import BigNumber from 'bignumber.js';
import { gasConfig, drtConfig } from 'src/config';
import { RewaOrDcdtTokenPayment } from 'src/models/dcdtTokenPayment.model';
import { decimalToHex } from 'src/utils/token.converters';
import { WrapAbiService } from 'src/modules/wrapping/services/wrap.abi.service';
import { TransactionOptions } from 'src/modules/common/transaction.options';

export type ComposableTask = {
    type: ComposableTaskType;
    arguments: BytesValue[];
};

@Injectable()
export class ComposableTasksTransactionService {
    constructor(
        private readonly drtProxy: MXProxyService,
        private readonly wrapAbi: WrapAbiService,
    ) {}

    async getComposeTasksTransaction(
        sender: string,
        payment: DcdtTokenPayment,
        tokenOut: RewaOrDcdtTokenPayment,
        tasks: ComposableTask[],
    ): Promise<TransactionModel> {
        let gasLimit: number = gasConfig.composableTasks.default;

        for (const task of tasks) {
            switch (task.type) {
                case ComposableTaskType.WRAP_REWA:
                    gasLimit += gasConfig.wrapREWA;
                    break;
                case ComposableTaskType.UNWRAP_REWA:
                    gasLimit += gasConfig.wrapREWA;
                    break;
                case ComposableTaskType.SWAP:
                    gasLimit +=
                        gasConfig.pairs.swapTokensFixedOutput.withFeeSwap;
                case ComposableTaskType.ROUTER_SWAP:
                    const routes = Math.trunc(task.arguments.length / 4);
                    gasLimit +=
                        routes * gasConfig.router.multiPairSwapMultiplier;
                default:
                    break;
            }
        }

        const transactionOptions = new TransactionOptions({
            sender: sender,
            chainID: drtConfig.chainID,
            gasLimit: gasLimit,
            function: 'composeTasks',
            arguments: [
                new Struct(RewaOrDcdtTokenPayment.getStructure(), [
                    new Field(
                        new TokenIdentifierValue(tokenOut.tokenIdentifier),
                        'token_identifier',
                    ),
                    new Field(new U64Value(new BigNumber(0)), 'token_nonce'),
                    new Field(
                        new BigUIntValue(new BigNumber(tokenOut.amount)),
                        'amount',
                    ),
                ]),
                VariadicValue.fromItems(...this.getRawTasks(tasks)),
            ],
        });

        if (payment.tokenIdentifier === drtConfig.REWAIdentifier) {
            transactionOptions.nativeTransferAmount = payment.amount;
        } else {
            transactionOptions.tokenTransfers = [
                new TokenTransfer({
                    token: new Token({
                        identifier: payment.tokenIdentifier,
                    }),
                    amount: BigInt(payment.amount),
                }),
            ];
        }

        return this.drtProxy.getComposableTasksContractTransaction(
            transactionOptions,
        );
    }

    async wrapRewaAndSwapTransaction(
        sender: string,
        value: string,
        tokenOutID: string,
        tokenOutAmountMin: string,
        swapEndpoint: string,
    ): Promise<TransactionModel> {
        const wrapTask: ComposableTask = {
            type: ComposableTaskType.WRAP_REWA,
            arguments: [],
        };

        const swapTask: ComposableTask = {
            type: ComposableTaskType.SWAP,
            arguments: [
                new BytesValue(Buffer.from(swapEndpoint, 'utf-8')),
                new BytesValue(Buffer.from(tokenOutID, 'utf-8')),
                new BytesValue(
                    Buffer.from(
                        decimalToHex(new BigNumber(tokenOutAmountMin)),
                        'hex',
                    ),
                ),
            ],
        };

        return this.getComposeTasksTransaction(
            sender,
            new DcdtTokenPayment({
                tokenIdentifier: 'REWA',
                tokenNonce: 0,
                amount: value,
            }),
            new RewaOrDcdtTokenPayment({
                tokenIdentifier: tokenOutID,
                amount: tokenOutAmountMin,
            }),
            [wrapTask, swapTask],
        );
    }

    async swapAndUnwrapRewaTransaction(
        sender: string,
        payment: DcdtTokenPayment,
        minimumValue: string,
        swapEndpoint: string,
    ): Promise<TransactionModel> {
        const wrappedRewaTokenID = await this.wrapAbi.wrappedRewaTokenID();

        const swapTask: ComposableTask = {
            type: ComposableTaskType.SWAP,
            arguments: [
                new BytesValue(Buffer.from(swapEndpoint, 'utf-8')),
                new BytesValue(Buffer.from(wrappedRewaTokenID, 'utf-8')),
                new BytesValue(
                    Buffer.from(
                        decimalToHex(new BigNumber(minimumValue)),
                        'hex',
                    ),
                ),
            ],
        };
        const unwrapTask: ComposableTask = {
            type: ComposableTaskType.UNWRAP_REWA,
            arguments: [],
        };

        return this.getComposeTasksTransaction(
            sender,
            payment,
            new RewaOrDcdtTokenPayment({
                tokenIdentifier: 'REWA',
                amount: minimumValue,
            }),
            [swapTask, unwrapTask],
        );
    }

    private getRawTasks(tasks: ComposableTask[]): TypedValue[] {
        const rawTasks: TypedValue[] = [];

        tasks.forEach((task) => {
            rawTasks.push(
                new EnumValue(
                    ComposableTaskEnumType.getEnumType(),
                    new EnumVariantDefinition(
                        task.type,
                        ComposableTaskEnumType.getEnumType().getVariantByName(
                            task.type,
                        ).discriminant,
                    ),
                    [],
                ),
            );
            rawTasks.push(
                new List(new ListType(new BytesType()), task.arguments),
            );
        });

        return rawTasks;
    }
}
