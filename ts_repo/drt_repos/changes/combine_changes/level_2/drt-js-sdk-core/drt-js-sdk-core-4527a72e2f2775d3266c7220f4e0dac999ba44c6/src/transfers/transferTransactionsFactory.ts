import { Address } from "../core/address";
import { REWA_IDENTIFIER_FOR_MULTI_DCDTNFT_TRANSFER } from "../core/constants";
import { ErrBadUsage } from "../core/errors";
import { TokenComputer, TokenTransfer } from "../core/tokens";
import { TokenTransfersDataBuilder } from "../core/tokenTransfersDataBuilder";
import { Transaction } from "../core/transaction";
import { TransactionBuilder } from "../core/transactionBuilder";
import * as resources from "./resources";

const ADDITIONAL_GAS_FOR_DCDT_TRANSFER = 100000;
const ADDITIONAL_GAS_FOR_DCDT_NFT_TRANSFER = 800000;

interface IConfig {
    chainID: string;
    minGasLimit: bigint;
    gasLimitPerByte: bigint;
    gasLimitDCDTTransfer: bigint;
    gasLimitDCDTNFTTransfer: bigint;
    gasLimitMultiDCDTNFTTransfer: bigint;
}

/**
 * Use this class to create transactions for native token transfers (REWA) or custom tokens transfers (DCDT/NTF/MetaDCDT).
 */
export class TransferTransactionsFactory {
    private readonly config?: IConfig;
    private readonly tokenTransfersDataBuilder?: TokenTransfersDataBuilder;
    private readonly tokenComputer?: TokenComputer;

    constructor(options: { config: IConfig }) {
        this.config = options.config;
        this.tokenComputer = new TokenComputer();
        this.tokenTransfersDataBuilder = new TokenTransfersDataBuilder();
    }

    createTransactionForNativeTokenTransfer(sender: Address, options: resources.NativeTokenTransferInput): Transaction {
        const data = options.data || new Uint8Array();

        return new Transaction({
            sender: sender,
            receiver: options.receiver,
            chainID: this.config!.chainID,
            gasLimit: this.computeGasForMoveBalance(this.config!, data),
            data: data,
            value: options.nativeAmount ?? BigInt(0),
        });
    }

    createTransactionForDCDTTokenTransfer(sender: Address, options: resources.CustomTokenTransferInput): Transaction {
        const numberOfTransfers = options.tokenTransfers.length;

        if (numberOfTransfers === 0) {
            throw new ErrBadUsage("No token transfer has been provided");
        }

        if (numberOfTransfers === 1) {
            return this.createSingleDCDTTransferTransaction(sender, options);
        }

        const { dataParts, extraGasForTransfer } = this.buildMultiDCDTNFTTransferData(
            options.tokenTransfers,
            options.receiver,
        );

        return new TransactionBuilder({
            config: this.config!,
            sender: sender,
            receiver: sender,
            dataParts: dataParts,
            gasLimit: extraGasForTransfer,
            addDataMovementGas: true,
        }).build();
    }

    createTransactionForTransfer(sender: Address, options: resources.CreateTransferTransactionInput): Transaction {
        const nativeAmount = options.nativeAmount ?? 0n;
        let tokenTransfers = options.tokenTransfers ? [...options.tokenTransfers] : [];
        const numberOfTokens = tokenTransfers.length;

        if (numberOfTokens && options.data?.length) {
            throw new ErrBadUsage("Can't set data field when sending dcdt tokens");
        }

        if ((nativeAmount && numberOfTokens === 0) || options.data) {
            return this.createTransactionForNativeTokenTransfer(sender, {
                receiver: options.receiver,
                nativeAmount: nativeAmount,
                data: options.data,
            });
        }

        const nativeTransfer = nativeAmount ? TokenTransfer.newFromNativeAmount(nativeAmount) : undefined;
        if (nativeTransfer) {
            tokenTransfers.push(nativeTransfer);
        }

        return this.createTransactionForDCDTTokenTransfer(sender, {
            receiver: options.receiver,
            tokenTransfers: tokenTransfers,
        });
    }

    private createSingleDCDTTransferTransaction(
        sender: Address,
        options: {
            receiver: Address;
            tokenTransfers: TokenTransfer[];
        },
    ): Transaction {
        const transfer = options.tokenTransfers[0];
        const { dataParts, extraGasForTransfer, receiver } = this.buildTransferData(transfer, {
            sender,
            receiver: options.receiver,
        });

        return new TransactionBuilder({
            config: this.config!,
            sender: sender,
            receiver: receiver,
            dataParts: dataParts,
            gasLimit: extraGasForTransfer,
            addDataMovementGas: true,
        }).build();
    }

    private buildTransferData(transfer: TokenTransfer, options: { sender: Address; receiver: Address }) {
        let dataParts: string[] = [];
        let extraGasForTransfer: bigint;
        let receiver = options.receiver;

        if (this.tokenComputer!.isFungible(transfer.token)) {
            if (transfer.token.identifier === REWA_IDENTIFIER_FOR_MULTI_DCDTNFT_TRANSFER) {
                ({ dataParts, extraGasForTransfer } = this.buildMultiDCDTNFTTransferData([transfer], receiver));
                receiver = options.sender;
            } else {
                ({ dataParts, extraGasForTransfer } = this.buildDCDTTransferData(transfer));
            }
        } else {
            ({ dataParts, extraGasForTransfer } = this.buildSingleDCDTNFTTransferData(transfer, receiver));
            receiver = options.sender; // Override receiver for non-fungible tokens
        }
        return { dataParts, extraGasForTransfer, receiver };
    }

    private buildMultiDCDTNFTTransferData(transfer: TokenTransfer[], receiver: Address) {
        return {
            dataParts: this.tokenTransfersDataBuilder!.buildDataPartsForMultiDCDTNFTTransfer(receiver, transfer),
            extraGasForTransfer:
                this.config!.gasLimitMultiDCDTNFTTransfer * BigInt(transfer.length) +
                BigInt(ADDITIONAL_GAS_FOR_DCDT_NFT_TRANSFER),
        };
    }

    private buildDCDTTransferData(transfer: TokenTransfer) {
        return {
            dataParts: this.tokenTransfersDataBuilder!.buildDataPartsForDCDTTransfer(transfer),
            extraGasForTransfer: this.config!.gasLimitDCDTTransfer + BigInt(ADDITIONAL_GAS_FOR_DCDT_TRANSFER),
        };
    }

    private buildSingleDCDTNFTTransferData(transfer: TokenTransfer, receiver: Address) {
        return {
            dataParts: this.tokenTransfersDataBuilder!.buildDataPartsForSingleDCDTNFTTransfer(transfer, receiver),
            extraGasForTransfer: this.config!.gasLimitDCDTNFTTransfer + BigInt(ADDITIONAL_GAS_FOR_DCDT_NFT_TRANSFER),
        };
    }

    private computeGasForMoveBalance(config: IConfig, data: Uint8Array): bigint {
        return config.minGasLimit + config.gasLimitPerByte * BigInt(data.length);
    }
}
