import {
    Address,
    AddressValue,
    BigUIntValue,
    BytesValue,
    Token,
    TokenTransfer,
    TypedValue,
    VariadicValue,
} from '@terradharitri/sdk-core';
import { Injectable } from '@nestjs/common';
import BigNumber from 'bignumber.js';
import { drtConfig, gasConfig } from 'src/config';
import { MultiSwapTokensArgs } from 'src/modules/auto-router/models/multi-swap-tokens.args';
import { WrapTransactionsService } from 'src/modules/wrapping/services/wrap.transactions.service';
import { TransactionModel } from '../../../models/transaction.model';
import { MXProxyService } from '../../../services/dharitri-communication/drt.proxy.service';
import { SWAP_TYPE } from '../models/auto-route.model';
import { ComposableTaskType } from 'src/modules/composable-tasks/models/composable.tasks.model';
import { ComposableTasksTransactionService } from 'src/modules/composable-tasks/services/composable.tasks.transaction';
import { DcdtTokenPayment } from '@terradharitri/sdk-exchange';
import { RewaOrDcdtTokenPayment } from 'src/models/dcdtTokenPayment.model';
import { decimalToHex } from 'src/utils/token.converters';
import { TransactionOptions } from 'src/modules/common/transaction.options';

@Injectable()
export class AutoRouterTransactionService {
    constructor(
        private readonly drtProxy: MXProxyService,
        private readonly transactionsWrapService: WrapTransactionsService,
        private readonly composeTasksTransactionService: ComposableTasksTransactionService,
    ) {}

    async multiPairSwap(
        sender: string,
        args: MultiSwapTokensArgs,
    ): Promise<TransactionModel[]> {
        const transactions = [];

        const amountIn = new BigNumber(args.intermediaryAmounts[0]).plus(
            new BigNumber(args.intermediaryAmounts[0]).multipliedBy(
                args.swapType === SWAP_TYPE.fixedOutput ? args.tolerance : 0,
            ),
        );

        if (args.tokenInID === drtConfig.REWAIdentifier) {
            return [
                await this.wrapRewaAndMultiSwapTransaction(
                    sender,
                    amountIn.integerValue().toFixed(),
                    args,
                ),
            ];
        }

        if (args.tokenOutID === drtConfig.REWAIdentifier) {
            return [
                await this.multiSwapAndUnwrapRewaTransaction(
                    sender,
                    amountIn.integerValue().toFixed(),
                    args,
                ),
            ];
        }

        const gasLimit =
            args.addressRoute.length * gasConfig.router.multiPairSwapMultiplier;

        const transactionOptions = new TransactionOptions({
            sender: sender,
            chainID: drtConfig.chainID,
            gasLimit: gasLimit,
            function: 'multiPairSwap',
            arguments:
                args.swapType == SWAP_TYPE.fixedInput
                    ? [
                          VariadicValue.fromItems(
                              ...this.multiPairFixedInputSwaps(args),
                          ),
                      ]
                    : [
                          VariadicValue.fromItems(
                              ...this.multiPairFixedOutputSwaps(args),
                          ),
                      ],
            tokenTransfers: [
                new TokenTransfer({
                    token: new Token({
                        identifier: args.tokenRoute[0],
                    }),
                    amount: BigInt(amountIn.integerValue().toFixed()),
                }),
            ],
        });

        const transaction =
            await this.drtProxy.getRouterSmartContractTransaction(
                transactionOptions,
            );
        transactions.push(transaction);

        if (args.tokenOutID === drtConfig.REWAIdentifier) {
            transactions.push(
                await this.transactionsWrapService.unwrapRewa(
                    sender,
                    args.intermediaryAmounts[
                        args.intermediaryAmounts.length - 1
                    ],
                ),
            );
        }

        return transactions;
    }

    multiPairFixedInputSwaps(args: MultiSwapTokensArgs): TypedValue[] {
        const swaps: TypedValue[] = [];

        const intermediaryTolerance = args.tolerance / args.addressRoute.length;

        for (const [index, address] of args.addressRoute.entries()) {
            const intermediaryToleranceMultiplier =
                args.addressRoute.length - index;

            const toleranceAmount = new BigNumber(
                args.intermediaryAmounts[index + 1],
            ).multipliedBy(
                intermediaryToleranceMultiplier * intermediaryTolerance,
            );

            const amountOutMin = new BigNumber(
                args.intermediaryAmounts[index + 1],
            )
                .minus(toleranceAmount)
                .integerValue();

            swaps.push(
                ...[
                    new AddressValue(Address.fromString(address)),
                    BytesValue.fromUTF8('swapTokensFixedInput'),
                    BytesValue.fromUTF8(args.tokenRoute[index + 1]),
                    new BigUIntValue(amountOutMin),
                ],
            );
        }
        return swaps;
    }

    multiPairFixedOutputSwaps(args: MultiSwapTokensArgs): TypedValue[] {
        const swaps: TypedValue[] = [];

        const intermediaryTolerance = args.tolerance / args.addressRoute.length;

        for (const [index, address] of args.addressRoute.entries()) {
            // method #1
            // [A -> B -> C -> D], all with swap_tokens_fixed_output
            // overall: less input, more gas, rest/dust in A, B & C
            const intermediaryToleranceMultiplier =
                args.addressRoute.length - index - 1;

            const toleranceAmount = new BigNumber(
                args.intermediaryAmounts[index + 1],
            ).multipliedBy(
                intermediaryToleranceMultiplier * intermediaryTolerance,
            );

            const amountOut = new BigNumber(args.intermediaryAmounts[index + 1])
                .plus(toleranceAmount)
                .integerValue()
                .toFixed();

            swaps.push(
                ...[
                    new AddressValue(Address.fromString(address)),
                    BytesValue.fromUTF8('swapTokensFixedOutput'),
                    BytesValue.fromUTF8(args.tokenRoute[index + 1]),
                    new BigUIntValue(new BigNumber(amountOut)),
                ],
            );
            // method #2
            // [A -> B -> C] with swap_tokens_fixed_input + [C -> D] with swap_tokens_fixed_output
            // overall: more input, less gas, rest in C
            /*if (index < args.addressRoute.length - 1) {
                const amountOutMin = new BigNumber(
                    args.intermediaryAmounts[index + 1],
                )
                    .plus(
                        new BigNumber(
                            args.intermediaryAmounts[index + 1],
                        ).multipliedBy(
                            (args.addressRoute.length - index - 1) *
                                intermediaryTolerance,
                        ),
                    )
                    .integerValue()
                    .toFixed();
                console.log(
                    'swapTokensFixedInput with amountOutMin ',
                    amountOutMin,
                );

                swaps.push(
                    ...[
                        new AddressValue(Address.fromString(address)),
                        BytesValue.fromUTF8('swapTokensFixedInput'),
                        BytesValue.fromUTF8(args.tokenRoute[index + 1]),
                        new BigUIntValue(
                            new BigNumber(args.intermediaryAmounts[index + 1]),
                        ),
                    ],
                );
            } else {
                console.log('swapTokensFixedOutput');
                swaps.push(
                    ...[
                        new AddressValue(Address.fromString(address)),
                        BytesValue.fromUTF8('swapTokensFixedOutput'),
                        BytesValue.fromUTF8(args.tokenRoute[index + 1]),
                        new BigUIntValue(
                            new BigNumber(args.intermediaryAmounts[index + 1]),
                        ),
                    ],
                );
            }*/
        }
        return swaps;
    }

    async wrapRewaAndMultiSwapTransaction(
        sender: string,
        value: string,
        args: MultiSwapTokensArgs,
    ): Promise<TransactionModel> {
        const typedArgs =
            args.swapType === SWAP_TYPE.fixedInput
                ? this.multiPairFixedInputSwaps(args)
                : this.multiPairFixedOutputSwaps(args);
        const swaps = this.convertMultiPairSwapsToBytesValues(typedArgs);

        return this.composeTasksTransactionService.getComposeTasksTransaction(
            sender,
            new DcdtTokenPayment({
                tokenIdentifier: 'REWA',
                tokenNonce: 0,
                amount: value,
            }),
            new RewaOrDcdtTokenPayment({
                tokenIdentifier: args.tokenRoute[args.tokenRoute.length - 1],
                nonce: 0,
                amount: args.intermediaryAmounts[
                    args.intermediaryAmounts.length - 1
                ],
            }),
            [
                {
                    type: ComposableTaskType.WRAP_REWA,
                    arguments: [],
                },
                {
                    type: ComposableTaskType.ROUTER_SWAP,
                    arguments: swaps,
                },
            ],
        );
    }

    async multiSwapAndUnwrapRewaTransaction(
        sender: string,
        value: string,
        args: MultiSwapTokensArgs,
    ): Promise<TransactionModel> {
        const typedArgs =
            args.swapType === SWAP_TYPE.fixedInput
                ? this.multiPairFixedInputSwaps(args)
                : this.multiPairFixedOutputSwaps(args);
        const swaps = this.convertMultiPairSwapsToBytesValues(typedArgs);

        return this.composeTasksTransactionService.getComposeTasksTransaction(
            sender,
            new DcdtTokenPayment({
                tokenIdentifier: args.tokenRoute[0],
                tokenNonce: 0,
                amount: value,
            }),
            new RewaOrDcdtTokenPayment({
                tokenIdentifier: 'REWA',
                nonce: 0,
                amount: args.intermediaryAmounts[
                    args.intermediaryAmounts.length - 1
                ],
            }),
            [
                {
                    type: ComposableTaskType.ROUTER_SWAP,
                    arguments: swaps,
                },
                {
                    type: ComposableTaskType.UNWRAP_REWA,
                    arguments: [],
                },
            ],
        );
    }

    private convertMultiPairSwapsToBytesValues(
        args: TypedValue[],
    ): BytesValue[] {
        if (args.length % 4 !== 0) {
            throw new Error('Invalid number of router swap arguments');
        }

        const swaps: BytesValue[] = [];

        for (let index = 0; index <= args.length - 4; index += 4) {
            const pairAddress = args[index];
            const functionName = args[index + 1];
            const tokenOutID = args[index + 2];
            const amountOutMin = args[index + 3];

            swaps.push(
                new BytesValue(Buffer.from(pairAddress.valueOf().hex(), 'hex')),
            );
            swaps.push(BytesValue.fromUTF8(functionName.valueOf()));
            swaps.push(BytesValue.fromUTF8(tokenOutID.valueOf()));
            swaps.push(
                new BytesValue(
                    Buffer.from(
                        decimalToHex(new BigNumber(amountOutMin.valueOf())),
                        'hex',
                    ),
                ),
            );
        }
        return swaps;
    }
}
