import { assert } from "chai";
import { Address, ErrBadUsage, Token, TokenTransfer } from "../core";
import { REWA_IDENTIFIER_FOR_MULTI_DCDTNFT_TRANSFER } from "../core/constants";
import { TransactionsFactoryConfig } from "../core/transactionsFactoryConfig";
import { TransferTransactionsFactory } from "./transferTransactionsFactory";

describe("test transfer transactions factory", function () {
    const config = new TransactionsFactoryConfig({ chainID: "D" });
    const transferFactory = new TransferTransactionsFactory({
        config: config,
    });

    const alice = Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
    const bob = Address.newFromBech32("drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c");

    it("should throw error, no token transfer provided", async () => {
        let transfers: any = [];

        assert.throw(
            () => {
                transferFactory.createTransactionForDCDTTokenTransfer(alice, {
                    receiver: bob,
                    tokenTransfers: transfers,
                });
            },
            ErrBadUsage,
            "No token transfer has been provided",
        );
    });

    it("should create 'Transaction' for native token transfer without data", async () => {
        const transaction = transferFactory.createTransactionForNativeTokenTransfer(alice, {
            receiver: bob,
            nativeAmount: 1000000000000000000n,
        });

        assert.equal(transaction.sender, alice);
        assert.equal(transaction.receiver, bob);
        assert.equal(transaction.value.valueOf(), 1000000000000000000n);
        assert.equal(transaction.gasLimit.valueOf(), 50000n);
        assert.deepEqual(transaction.data, new Uint8Array());
    });

    it("should create 'Transaction' for native token transfer with data", async () => {
        const transaction = transferFactory.createTransactionForNativeTokenTransfer(alice, {
            receiver: bob,
            nativeAmount: 1000000000000000000n,
            data: Buffer.from("test data"),
        });

        assert.equal(transaction.sender, alice);
        assert.equal(transaction.receiver, bob);
        assert.equal(transaction.value.valueOf(), 1000000000000000000n);
        assert.equal(transaction.gasLimit.valueOf(), 63500n);
        assert.deepEqual(transaction.data, Buffer.from("test data"));
    });

    it("should create 'Transaction' for dcdt transfer", async () => {
        const fooToken = new Token({ identifier: "FOO-123456", nonce: 0n });
        const transfer = new TokenTransfer({ token: fooToken, amount: 1000000n });

        const transaction = transferFactory.createTransactionForDCDTTokenTransfer(alice, {
            receiver: bob,
            tokenTransfers: [transfer],
        });

        assert.equal(transaction.sender, alice);
        assert.equal(transaction.receiver, bob);
        assert.equal(transaction.value.valueOf(), 0n);
        assert.equal(transaction.gasLimit.valueOf(), 410000n);
        assert.deepEqual(transaction.data.toString(), "DCDTTransfer@464f4f2d313233343536@0f4240");
    });

    it("should create 'Transaction' for nft transfer", async () => {
        const nft = new Token({ identifier: "NFT-123456", nonce: 10n });
        const transfer = new TokenTransfer({ token: nft, amount: 1n });

        const transaction = transferFactory.createTransactionForDCDTTokenTransfer(alice, {
            receiver: bob,
            tokenTransfers: [transfer],
        });

        assert.equal(transaction.sender, alice);
        assert.equal(transaction.receiver, alice);
        assert.equal(transaction.value.valueOf(), 0n);
        assert.equal(transaction.gasLimit.valueOf(), 1210500n);
        assert.deepEqual(
            transaction.data.toString(),
            "DCDTNFTTransfer@4e46542d313233343536@0a@01@8049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f8",
        );
    });

    it("should create 'Transaction' for nft transfer with prefix", async () => {
        const nft = new Token({ identifier: "t0-NFT-123456", nonce: 10n });
        const transfer = new TokenTransfer({ token: nft, amount: 1n });

        const transaction = transferFactory.createTransactionForDCDTTokenTransfer(alice, {
            receiver: bob,
            tokenTransfers: [transfer],
        });

        assert.equal(transaction.sender, alice);
        assert.equal(transaction.receiver, alice);
        assert.equal(transaction.value.valueOf(), 0n);
        assert.equal(transaction.gasLimit.valueOf(), 1219500n);
        assert.deepEqual(
            transaction.data.toString(),
            "DCDTNFTTransfer@74302d4e46542d313233343536@0a@01@8049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f8",
        );
    });

    it("should create 'Transaction' for multiple nft transfers", async () => {
        const firstNft = new Token({ identifier: "NFT-123456", nonce: 10n });
        const firstTransfer = new TokenTransfer({ token: firstNft, amount: 1n });

        const secondNft = new Token({ identifier: "TEST-987654", nonce: 1n });
        const secondTransfer = new TokenTransfer({ token: secondNft, amount: 1n });

        const transaction = transferFactory.createTransactionForDCDTTokenTransfer(alice, {
            receiver: bob,
            tokenTransfers: [firstTransfer, secondTransfer],
        });

        assert.equal(transaction.sender, alice);
        assert.equal(transaction.receiver, alice);
        assert.equal(transaction.value.valueOf(), 0n);
        assert.equal(transaction.gasLimit.valueOf(), 1466000n);
        assert.deepEqual(
            transaction.data.toString(),
            "MultiDCDTNFTTransfer@8049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f8@02@4e46542d313233343536@0a@01@544553542d393837363534@01@01",
        );

        const secondTransaction = transferFactory.createTransactionForTransfer(alice, {
            receiver: bob,
            tokenTransfers: [firstTransfer, secondTransfer],
        });

        assert.deepEqual(transaction, secondTransaction);
    });

    it("should create 'Transaction' for multiple nft transfers with prefix", async () => {
        const firstNft = new Token({ identifier: "t0-NFT-123456", nonce: 10n });
        const firstTransfer = new TokenTransfer({ token: firstNft, amount: 1n });

        const secondNft = new Token({ identifier: "t0-TEST-987654", nonce: 1n });
        const secondTransfer = new TokenTransfer({ token: secondNft, amount: 1n });

        const transaction = transferFactory.createTransactionForDCDTTokenTransfer(alice, {
            receiver: bob,
            tokenTransfers: [firstTransfer, secondTransfer],
        });

        assert.equal(transaction.sender, alice);
        assert.equal(transaction.receiver, alice);
        assert.equal(transaction.value.valueOf(), 0n);
        assert.equal(transaction.gasLimit.valueOf(), 1484000n);
        assert.deepEqual(
            transaction.data.toString(),
            "MultiDCDTNFTTransfer@8049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f8@02@74302d4e46542d313233343536@0a@01@74302d544553542d393837363534@01@01",
        );

        const secondTransaction = transferFactory.createTransactionForTransfer(alice, {
            receiver: bob,
            tokenTransfers: [firstTransfer, secondTransfer],
        });

        assert.deepEqual(transaction, secondTransaction);
    });

    it("should fail to create transaction for token transfers", async () => {
        assert.throws(() => {
            const nft = new Token({ identifier: "NFT-123456", nonce: 10n });
            const transfer = new TokenTransfer({ token: nft, amount: 1n });

            transferFactory.createTransactionForTransfer(alice, {
                receiver: bob,
                tokenTransfers: [transfer],
                data: Buffer.from("test"),
            });
        }, "Can't set data field when sending dcdt tokens");
    });

    it("should create transaction for native transfers", async () => {
        const transaction = transferFactory.createTransactionForTransfer(alice, {
            receiver: bob,
            nativeAmount: 1000000000000000000n,
        });

        assert.equal(transaction.sender, alice);
        assert.equal(transaction.receiver, bob);
        assert.equal(transaction.value.valueOf(), 1000000000000000000n);
        assert.equal(transaction.gasLimit.valueOf(), 50000n);
    });

    it("should create transaction for native transfers and set data field", async () => {
        const transaction = transferFactory.createTransactionForTransfer(alice, {
            receiver: bob,
            nativeAmount: 1000000000000000000n,
            data: Buffer.from("hello"),
        });

        assert.equal(transaction.sender, alice);
        assert.equal(transaction.receiver, bob);
        assert.equal(transaction.value.valueOf(), 1000000000000000000n);
        assert.equal(transaction.gasLimit.valueOf(), 57500n);
        assert.deepEqual(transaction.data, Buffer.from("hello"));
    });

    it("should create transaction for notarizing", async () => {
        const transaction = transferFactory.createTransactionForTransfer(alice, {
            receiver: bob,
            data: Buffer.from("hello"),
        });

        assert.equal(transaction.sender, alice);
        assert.equal(transaction.receiver, bob);
        assert.equal(transaction.gasLimit.valueOf(), 57500n);
        assert.deepEqual(transaction.data, Buffer.from("hello"));
    });

    it("should create transaction for token transfers", async () => {
        const firstNft = new Token({ identifier: "NFT-123456", nonce: 10n });
        const firstTransfer = new TokenTransfer({ token: firstNft, amount: 1n });

        const secondNft = new Token({ identifier: "TEST-987654", nonce: 1n });
        const secondTransfer = new TokenTransfer({ token: secondNft, amount: 1n });

        const transaction = transferFactory.createTransactionForTransfer(alice, {
            receiver: bob,
            nativeAmount: 1000000000000000000n,
            tokenTransfers: [firstTransfer, secondTransfer],
        });

        assert.equal(transaction.sender, alice);
        assert.equal(transaction.receiver, alice);
        assert.equal(transaction.value.valueOf(), 0n);
        assert.equal(transaction.gasLimit.valueOf(), 1727500n);
        assert.deepEqual(
            transaction.data.toString(),
            "MultiDCDTNFTTransfer@8049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f8@03@4e46542d313233343536@0a@01@544553542d393837363534@01@01@524557412d303030303030@@0de0b6b3a7640000",
        );
    });

    it("should create transaction for token transfers with prefix", async () => {
        const firstNft = new Token({ identifier: "t0-NFT-123456", nonce: 10n });
        const firstTransfer = new TokenTransfer({ token: firstNft, amount: 1n });

        const secondNft = new Token({ identifier: "t0-TEST-987654", nonce: 1n });
        const secondTransfer = new TokenTransfer({ token: secondNft, amount: 1n });

        const transaction = transferFactory.createTransactionForTransfer(alice, {
            receiver: bob,
            nativeAmount: 1000000000000000000n,
            tokenTransfers: [firstTransfer, secondTransfer],
        });

        assert.equal(transaction.sender, alice);
        assert.equal(transaction.receiver, alice);
        assert.equal(transaction.value.valueOf(), 0n);
        assert.equal(transaction.gasLimit.valueOf(), 1745500n);
        assert.deepEqual(
            transaction.data.toString(),
            "MultiDCDTNFTTransfer@8049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f8@03@74302d4e46542d313233343536@0a@01@74302d544553542d393837363534@01@01@524557412d303030303030@@0de0b6b3a7640000",
        );
    });

    it("should create multi transfer for rewa", async () => {
        const firstNft = new Token({ identifier: REWA_IDENTIFIER_FOR_MULTI_DCDTNFT_TRANSFER });
        const firstTransfer = new TokenTransfer({ token: firstNft, amount: 1000000000000000000n });

        const transaction = transferFactory.createTransactionForDCDTTokenTransfer(alice, {
            receiver: bob,
            tokenTransfers: [firstTransfer],
        });

        assert.equal(transaction.sender, alice);
        assert.equal(transaction.receiver, alice);
        assert.equal(transaction.value.valueOf(), 0n);
        assert.equal(transaction.gasLimit.valueOf(), 1_243_500n);
        assert.deepEqual(
            transaction.data.toString(),
            "MultiDCDTNFTTransfer@8049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f8@01@524557412d303030303030@@0de0b6b3a7640000",
        );
    });
});
