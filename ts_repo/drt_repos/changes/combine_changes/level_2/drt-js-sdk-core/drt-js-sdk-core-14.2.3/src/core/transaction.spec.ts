import { Buffer } from "buffer";
import { assert } from "chai";
import { Account } from "../accounts";
import { ProtoSerializer } from "../proto";
import { getTestWalletsPath } from "../testutils/utils";
import { Address } from "./address";
import { MIN_TRANSACTION_VERSION_THAT_SUPPORTS_OPTIONS, TRANSACTION_OPTIONS_DEFAULT } from "./constants";
import { INetworkConfig } from "./interfaces";
import { Transaction } from "./transaction";
import { TransactionComputer } from "./transactionComputer";

describe("test transaction", async () => {
    let alice: Account;
    let bob: Account;
    let carol: Account;
    const minGasLimit = 50000;
    const minGasPrice = 1000000000;

    const transactionComputer = new TransactionComputer();

    const networkConfig: INetworkConfig = {
        minGasLimit: 50000n,
        gasPerDataByte: 1500n,
        gasPriceModifier: 0.01,
        chainID: "D",
    };

    before(async function () {
        alice = await Account.newFromPem(`${getTestWalletsPath()}/alice.pem`);
        bob = await Account.newFromPem(`${getTestWalletsPath()}/bob.pem`);
        carol = await Account.newFromPem(`${getTestWalletsPath()}/carol.pem`);
    });

    it("should serialize transaction for signing (without data)", async () => {
        const transaction = new Transaction({
            chainID: networkConfig.chainID,
            sender: alice.address,
            receiver: bob.address,
            gasLimit: 50000n,
            value: 0n,
            version: 2,
            nonce: 89n,
        });

        const serializedTransactionBytes = transactionComputer.computeBytesForSigning(transaction);
        const serializedTransaction = Buffer.from(serializedTransactionBytes).toString();

        assert.equal(
            serializedTransaction,
            `{"nonce":89,"value":"0","receiver":"drt18h03w0y7qtqwtra3u4f0gu7e3kn2fslj83lqxny39m5c4rwaectswerhd2","sender":"drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh","gasPrice":1000000000,"gasLimit":50000,"chainID":"D","version":2}`,
        );
    });

    it("should serialize transaction for signing (with data)", async () => {
        const transaction = new Transaction({
            chainID: networkConfig.chainID,
            sender: alice.address,
            receiver: bob.address,
            gasLimit: 70000n,
            value: 1000000000000000000n,
            version: 2,
            nonce: 90n,
            data: new Uint8Array(Buffer.from("hello")),
        });

        const serializedTransactionBytes = transactionComputer.computeBytesForSigning(transaction);
        const serializedTransaction = Buffer.from(serializedTransactionBytes).toString();

        assert.equal(
            serializedTransaction,
            `{"nonce":90,"value":"1000000000000000000","receiver":"drt18h03w0y7qtqwtra3u4f0gu7e3kn2fslj83lqxny39m5c4rwaectswerhd2","sender":"drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh","gasPrice":1000000000,"gasLimit":70000,"data":"aGVsbG8=","chainID":"D","version":2}`,
        );
    });

    it("should sign transaction (with no data, no value) (legacy)", async () => {
        const transaction = new Transaction({
            nonce: 89n,
            value: 0n,
            sender: alice.address,
            receiver: bob.address,
            gasPrice: BigInt(minGasPrice),
            gasLimit: BigInt(minGasLimit),
            chainID: "local-testnet",
        });

        transaction.signature = await alice.signTransaction(transaction);

        assert.equal(
            Buffer.from(transaction.signature).toString("hex"),
            "129faad6b4222905addd10a2b49941f0a0f559de52e78af98752916cd470980527f4822a8503ded124f6a54f2801c2617c1922a6d0d2c35c33daa437311a4903",
        );
        assert.equal(
            transactionComputer.computeTransactionHash(transaction),
            "2014e16bbfd0964fbcb37e4d76eef76cac479ce295109063e28535b66360f6ea",
        );
    });

    it("should sign transaction (with data, no value) (legacy)", async () => {
        const transaction = new Transaction({
            nonce: 90n,
            value: 0n,
            sender: alice.address,
            receiver: bob.address,
            gasPrice: BigInt(minGasPrice),
            gasLimit: 80000n,
            data: Buffer.from("hello"),
            chainID: "local-testnet",
        });

        transaction.signature = await alice.signTransaction(transaction);

        assert.equal(
            Buffer.from(transaction.signature).toString("hex"),
            "32290b79b021945a6650c5f38bcc8bf70f025323598396eaae68447912a84b5bc07897cf5261c48de1b0ab588148811cc394a9462c65372c8dfed1e245951807",
        );
        assert.equal(
            transactionComputer.computeTransactionHash(transaction),
            "641c2107f1035be443697fd18ba95a19197ce54b2e85d64e814de12fb30ae2ba",
        );
    });

    it("should sign transaction (with usernames)", async () => {
        const transaction = new Transaction({
            chainID: "T",
            sender: carol.address,
            receiver: alice.address,
            gasLimit: 50000n,
            value: 1000000000000000000n,
            version: 2,
            nonce: 204n,
            senderUsername: "carol",
            receiverUsername: "alice",
        });

        transaction.signature = await carol.signTransaction(transaction);

        assert.equal(
            Buffer.from(transaction.signature).toString("hex"),
            "6e217efb1106be3fbc37065b16491112cb7f1adc78b8860f6d587733776afdce0c9ce7f36f102c5b1fe14a2a4a7844812fa0edb98451deecc8b0225f5250340c",
        );
    });

    it("should compute hash", async () => {
        const transaction = new Transaction({
            chainID: networkConfig.chainID,
            sender: alice.address,
            receiver: alice.address,
            gasLimit: 100000n,
            value: 1000000000000n,
            version: 2,
            nonce: 17243n,
            data: Buffer.from("testtx"),
        });

        transaction.signature = Buffer.from(
            "eaa9e4dfbd21695d9511e9754bde13e90c5cfb21748a339a79be11f744c71872e9fe8e73c6035c413f5f08eef09e5458e9ea6fc315ff4da0ab6d000b450b2a07",
            "hex",
        );

        const hash = transactionComputer.computeTransactionHash(transaction);

        assert.equal(hash, "60a8b538c87c517cff5fd2ca65d9eb4a1628ed7812a41e7d0f48b43c420cc79d");
    });

    it("should compute hash (with usernames)", async () => {
        const transaction = new Transaction({
            chainID: networkConfig.chainID,
            sender: alice.address,
            receiver: alice.address,
            gasLimit: 100000n,
            value: 1000000000000n,
            version: 2,
            nonce: 17244n,
            data: Buffer.from("testtx"),
            senderUsername: "alice",
            receiverUsername: "alice",
        });

        transaction.signature = Buffer.from(
            "807bcd7de5553ea6dfc57c0510e84d46813c5963d90fec50991c500091408fcf6216dca48dae16a579a1611ed8b2834bae8bd0027dc17eb557963f7151b82c07",
            "hex",
        );

        const hash = transactionComputer.computeTransactionHash(transaction);

        assert.equal(hash, "ad24bbd78111eb9cc4f3c7ec90dfd9fe2e9b76fb7af8dd8839cda6605ece6591");
    });

    it("should sign & compute hash (with data, with opaque, unused options) (legacy)", async () => {
        const transaction = new Transaction({
            nonce: 89n,
            value: 0n,
            sender: alice.address,
            receiver: bob.address,
            gasPrice: BigInt(minGasPrice),
            gasLimit: BigInt(minGasLimit),
            chainID: "local-testnet",
            // The protocol ignores the options when version == 1
            version: 1,
            options: 1,
        });

        assert.throws(() => {
            transactionComputer.computeBytesForSigning(transaction);
        }, `Non-empty transaction options requires transaction version >= ${MIN_TRANSACTION_VERSION_THAT_SUPPORTS_OPTIONS}`);
    });

    it("should sign & compute hash (with data, with value) (legacy)", async () => {
        const transaction = new Transaction({
            nonce: 91n,
            value: 10000000000000000000n,
            sender: alice.address,
            receiver: bob.address,
            gasPrice: BigInt(minGasPrice),
            gasLimit: 100000n,
            data: Buffer.from("for the book"),
            chainID: "local-testnet",
        });

        transaction.signature = await alice.signTransaction(transaction);

        assert.equal(
            Buffer.from(transaction.signature).toString("hex"),
            "2605d0c7298a218608fe5324cb9bdea34f52b301b73dc8a351c902220e880c50ad2b4039e68ac403103b24eaa4e8bc32fbfa90a58da6dce818b8491929f6550d",
        );
        assert.equal(
            transactionComputer.computeTransactionHash(transaction).toString(),
            "e96e41cb4b78c72c1a2241a611c1ffd3edd76804bf82470c5111be7cc6192aba",
        );
    });

    it("should sign & compute hash (with data, with large value) (legacy)", async () => {
        const transaction = new Transaction({
            nonce: 92n,
            value: BigInt("123456789000000000000000000000"),
            sender: alice.address,
            receiver: bob.address,
            gasPrice: BigInt(minGasPrice),
            gasLimit: 100000n,
            data: Buffer.from("for the spaceship"),
            chainID: "local-testnet",
        });

        transaction.signature = await alice.signTransaction(transaction);

        assert.equal(
            Buffer.from(transaction.signature).toString("hex"),
            "61deaf3566fc651e5d50ae8f8989771c149077b8f721e02e886be66f013be4e4ae1eadc6d36d54ce0d3899b34db6e6198d01f7b155243305fb49199f21937609",
        );
        assert.equal(
            transactionComputer.computeTransactionHash(transaction),
            "429a7f3475030ef1e15fc326336c20fa3d67e01f4901718c5c7c1be3cd8a7c7e",
        );
    });

    it("should sign & compute hash (with nonce = 0) (legacy)", async () => {
        const transaction = new Transaction({
            nonce: 0n,
            value: 0n,
            sender: alice.address,
            receiver: bob.address,
            gasPrice: BigInt(minGasPrice),
            gasLimit: 80000n,
            data: Buffer.from("hello"),
            chainID: "local-testnet",
            version: 1,
        });

        transaction.signature = await alice.signTransaction(transaction);

        assert.equal(
            Buffer.from(transaction.signature).toString("hex"),
            "24b30a91becdbba0dddc74886d52253a221ed2c381da0c6d7ea4c6482075b577c1b105c588ba8ec9442182b345e3cda975ecfc3785481ad44eb2dfd0464aee00",
        );
        assert.equal(
            transactionComputer.computeTransactionHash(transaction),
            "f9b37683eefa8397a78232268e57b85ea2989145bd2f9ecbc2a366c9f50051db",
        );
    });

    it("should sign & compute hash (without options field, should be omitted) (legacy)", async () => {
        const transaction = new Transaction({
            nonce: 89n,
            value: 0n,
            sender: alice.address,
            receiver: bob.address,
            gasPrice: BigInt(minGasPrice),
            gasLimit: BigInt(minGasLimit),
            chainID: "local-testnet",
        });

        transaction.signature = await alice.signTransaction(transaction);

        assert.equal(
            Buffer.from(transaction.signature).toString("hex"),
            "129faad6b4222905addd10a2b49941f0a0f559de52e78af98752916cd470980527f4822a8503ded124f6a54f2801c2617c1922a6d0d2c35c33daa437311a4903",
        );
        assert.equal(
            transactionComputer.computeTransactionHash(transaction),
            "2014e16bbfd0964fbcb37e4d76eef76cac479ce295109063e28535b66360f6ea",
        );

        const result = transactionComputer.computeBytesForSigning(transaction);
        assert.isFalse(result.toString().includes("options"));
    });

    it("should sign & compute hash (with guardian field, should be omitted) (legacy)", async () => {
        const transaction = new Transaction({
            nonce: 89n,
            value: 0n,
            sender: alice.address,
            receiver: bob.address,
            gasPrice: BigInt(minGasPrice),
            gasLimit: BigInt(minGasLimit),
            chainID: "local-testnet",
        });

        transaction.signature = await alice.signTransaction(transaction);

        assert.equal(
            Buffer.from(transaction.signature).toString("hex"),
            "129faad6b4222905addd10a2b49941f0a0f559de52e78af98752916cd470980527f4822a8503ded124f6a54f2801c2617c1922a6d0d2c35c33daa437311a4903",
        );
        assert.equal(
            transactionComputer.computeTransactionHash(transaction),
            "2014e16bbfd0964fbcb37e4d76eef76cac479ce295109063e28535b66360f6ea",
        );

        const result = transactionComputer.computeBytesForSigning(transaction);
        assert.isFalse(result.toString().includes("options"));
    });

    it("should sign & compute hash (with usernames) (legacy)", async () => {
        const transaction = new Transaction({
            nonce: 204n,
            value: 1000000000000000000n,
            sender: Address.newFromBech32("drt1k2s324ww2g0yj38qn2ch2jwctdy8mnfxep94q9arncc6xecg3xaq889n6e"),
            receiver: Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            senderUsername: "carol",
            receiverUsername: "alice",
            gasLimit: 50000n,
            chainID: "T",
        });

        transaction.signature = await carol.signTransaction(transaction);

        assert.equal(
            Buffer.from(transaction.signature).toString("hex"),
            "c6aee90d57933382f0b06ce75619b7a4dfaa05f22f885dee7272654ea6bf96b66f264ab8fbe84282a0817893c046e9f2dd875621655b3215b4fc6ba251ccd808",
        );
        assert.equal(
            transactionComputer.computeTransactionHash(transaction),
            "d2015b71186707cdf8bfe6a005c6fe8c5d4fafa02b7ce6743042d1a54302796c",
        );
    });

    it("should sign & compute hash (guarded transaction)", async () => {
        const transaction = new Transaction({
            chainID: "local-testnet",
            sender: alice.address,
            receiver: bob.address,
            gasLimit: 150000n,
            gasPrice: 1000000000n,
            data: new Uint8Array(Buffer.from("test data field")),
            version: 2,
            options: 2,
            nonce: 92n,
            value: 123456789000000000000000000000n,
            guardian: Address.newFromBech32("drt1x23lzn8483xs2su4fak0r0dqx6w38enpmmqf2yrkylwq7mfnvyhsmueha6"),
        });
        transaction.guardianSignature = new Uint8Array(64);
        transaction.signature = await alice.signTransaction(transaction);

        const serializer = new ProtoSerializer();
        const buffer = serializer.serializeTransaction(transaction);

        assert.equal(
            buffer.toString("hex"),
            "085c120e00018ee90ff6181f3761632000001a203ddf173c9e02c0e58fb1e552f473d98da6a4c3f23c7e034c912ee98a8dddce172a20391f932707a9dfa86d3bcbb3d5d0cc9f25ad0e680fe499f107d844b7e6ea71d5388094ebdc0340f093094a0f746573742064617461206669656c64520d6c6f63616c2d746573746e65745802624014bbdaf3ece1533aefe874147266fd8f6d7281a571a12d133c7dfb3cda655a8618c1092bc22c4496eaf867c3af3686fc4ae7327485a13c12769fcd589dbc2a0d6802722032a3f14cf53c4d0543954f6cf1bda0369d13e661dec095107627dc0f6d33612f7a4000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
        );

        const txHash = transactionComputer.computeTransactionHash(transaction);
        assert.equal(txHash, "5b1a94420f59e687200044233abbddb3056e6921deb82cec706ecbd03003fb93");
    });

    it("computes fee (legacy)", () => {
        const transaction = new Transaction({
            nonce: 92n,
            value: BigInt("123456789000000000000000000000"),
            sender: alice.address,
            receiver: bob.address,
            gasPrice: BigInt(minGasPrice),
            gasLimit: BigInt(minGasLimit),
            chainID: "local-testnet",
        });

        const fee = transactionComputer.computeTransactionFee(transaction, networkConfig);
        assert.equal(fee.toString(), "50000000000000");
    });

    it("computes fee", async () => {
        const transaction = new Transaction({
            chainID: "D",
            sender: alice.address,
            receiver: alice.address,
            gasLimit: 50000n,
            gasPrice: BigInt(minGasPrice),
        });

        const gasLimit = transactionComputer.computeTransactionFee(transaction, networkConfig);
        assert.equal(gasLimit.toString(), "50000000000000");
    });

    it("computes fee, but should throw `NotEnoughGas` error", async () => {
        const transaction = new Transaction({
            chainID: networkConfig.chainID,
            sender: alice.address,
            receiver: alice.address,
            gasLimit: 50000n,
            data: Buffer.from("toolittlegaslimit"),
        });

        assert.throws(() => {
            transactionComputer.computeTransactionFee(transaction, networkConfig);
        });
    });

    it("computes fee (with data field) (legacy)", () => {
        let transaction = new Transaction({
            nonce: 92n,
            value: BigInt("123456789000000000000000000000"),
            sender: alice.address,
            receiver: bob.address,
            data: Buffer.from("testdata"),
            gasPrice: BigInt(minGasPrice),
            gasLimit: BigInt(minGasLimit + 12010),
            chainID: "local-testnet",
        });

        let fee = transactionComputer.computeTransactionFee(transaction, networkConfig);
        assert.equal(fee.toString(), "62000100000000");
    });

    it("computes fee (with data field)", async () => {
        const transaction = new Transaction({
            chainID: networkConfig.chainID,
            sender: alice.address,
            receiver: alice.address,
            gasLimit: 50000n + 12010n,
            gasPrice: BigInt(minGasPrice),
            data: Buffer.from("testdata"),
        });

        const gasLimit = transactionComputer.computeTransactionFee(transaction, networkConfig);
        assert.equal(gasLimit.toString(), "62000100000000");
    });

    it("should convert transaction to plain object and back", () => {
        const sender = alice.address;
        const transaction = new Transaction({
            nonce: 90n,
            value: 123456789000000000000000000000n,
            sender: sender,
            receiver: bob.address,
            senderUsername: "alice",
            receiverUsername: "bob",
            gasPrice: BigInt(minGasPrice),
            gasLimit: 80000n,
            data: Buffer.from("hello"),
            chainID: "local-testnet",
        });

        const plainObject = transaction.toPlainObject();
        const restoredTransaction = Transaction.newFromPlainObject(plainObject);
        assert.deepEqual(restoredTransaction, transaction);
    });

    it("should handle large values", () => {
        const tx1 = new Transaction({
            value: 123456789000000000000000000000n,
            sender: alice.address,
            receiver: bob.address,
            gasLimit: 50000n,
            chainID: "local-testnet",
        });
        assert.equal(tx1.value.toString(), "123456789000000000000000000000");

        const tx2 = new Transaction({
            value: 123456789000000000000000000000n,
            sender: alice.address,
            receiver: bob.address,
            gasLimit: 50000n,
            chainID: "local-testnet",
        });
        assert.equal(tx2.value.toString(), "123456789000000000000000000000");

        const tx3 = new Transaction({
            value: BigInt("123456789000000000000000000000"),
            sender: alice.address,
            receiver: bob.address,
            gasLimit: 50000n,
            chainID: "local-testnet",
        });
        assert.equal(tx3.value.toString(), "123456789000000000000000000000");
    });

    it("checks correctly the version and options of the transaction", async () => {
        let transaction = new Transaction({
            nonce: 90n,
            value: BigInt("1000000000000000000"),
            sender: alice.address,
            receiver: bob.address,
            gasPrice: BigInt(minGasPrice),
            gasLimit: 80000n,
            data: Buffer.from("hello"),
            chainID: "local-testnet",
            version: 1,
            options: TRANSACTION_OPTIONS_DEFAULT,
        });
        assert.isFalse(transaction.isGuardedTransaction());

        transaction = new Transaction({
            nonce: 90n,
            value: BigInt("1000000000000000000"),
            sender: alice.address,
            receiver: bob.address,
            gasPrice: BigInt(minGasPrice),
            gasLimit: 80000n,
            data: Buffer.from("hello"),
            chainID: "local-testnet",
            version: 1,
            options: 2,
        });
        assert.isFalse(transaction.isGuardedTransaction());

        transaction = new Transaction({
            nonce: 90n,
            value: BigInt("1000000000000000000"),
            sender: alice.address,
            receiver: bob.address,
            gasPrice: BigInt(minGasPrice),
            gasLimit: 80000n,
            data: Buffer.from("hello"),
            chainID: "local-testnet",
            version: 2,
            options: 2,
        });
        assert.isFalse(transaction.isGuardedTransaction());

        transaction = new Transaction({
            nonce: 90n,
            value: BigInt("1000000000000000000"),
            sender: alice.address,
            receiver: bob.address,
            gasPrice: BigInt(minGasPrice),
            gasLimit: 80000n,
            data: Buffer.from("hello"),
            chainID: "local-testnet",
            version: 2,
            options: 2,
        });
        assert.isFalse(transaction.isGuardedTransaction());

        transaction = new Transaction({
            nonce: 90n,
            value: BigInt("1000000000000000000"),
            sender: alice.address,
            receiver: bob.address,
            gasPrice: BigInt(minGasPrice),
            guardian: bob.address,
            gasLimit: 80000n,
            data: Buffer.from("hello"),
            chainID: "local-testnet",
            version: 2,
            options: 2,
        });
        assert.isFalse(transaction.isGuardedTransaction());

        transaction = new Transaction({
            nonce: 90n,
            value: BigInt("1000000000000000000"),
            sender: alice.address,
            receiver: bob.address,
            gasPrice: BigInt(minGasPrice),
            guardian: bob.address,
            gasLimit: 80000n,
            data: Buffer.from("hello"),
            chainID: "local-testnet",
            version: 2,
            options: 2,
        });
        transaction.signature = await alice.signTransaction(transaction);
        transaction.guardianSignature = transaction.signature;
        assert.isTrue(transaction.isGuardedTransaction());
    });

    it("test sign using hash", async () => {
        let transaction = new Transaction({
            nonce: 89n,
            value: 0n,
            sender: alice.address,
            receiver: bob.address,
            gasLimit: 50000n,
            gasPrice: 1000000000n,
            chainID: "integration tests chain ID",
            version: 2,
            options: 1,
        });

        transaction.signature = await alice.sign(transactionComputer.computeHashForSigning(transaction));

        assert.equal(
            Buffer.from(transaction.signature).toString("hex"),
            "17e608f5ae9897b500046cf896cc5ea3de5208e8b42781733a9ae5a0f1ef11a5a0710b30b9df2bcc550fd83a379ffbb5dab20a6d803c755ca969a294f563ac08",
        );
    });

    it("should apply guardian", async () => {
        let transaction = new Transaction({
            nonce: 89n,
            value: 0n,
            sender: alice.address,
            receiver: bob.address,
            gasLimit: 50000n,
            chainID: "localnet",
        });

        transactionComputer.applyGuardian(transaction, carol.address);

        assert.equal(transaction.version, 2);
        assert.equal(transaction.options, 2);
        assert.equal(transaction.guardian, carol.address);
    });

    it("should apply guardian with options set for hash signing", async () => {
        let transaction = new Transaction({
            nonce: 89n,
            value: 0n,
            sender: alice.address,
            receiver: bob.address,
            gasLimit: 50000n,
            chainID: "localnet",
            version: 1,
        });

        transactionComputer.applyOptionsForHashSigning(transaction);
        assert.equal(transaction.version, 2);
        assert.equal(transaction.options, 1);

        transactionComputer.applyGuardian(transaction, carol.address);
        assert.equal(transaction.version, 2);
        assert.equal(transaction.options, 3);
    });

    it("should ensure transaction is valid", async () => {
        let transaction = new Transaction({
            sender: Address.empty(),
            receiver: bob.address,
            gasLimit: 50000n,
            chainID: "",
        });

        transaction.sender = alice.address;

        assert.throws(() => {
            transactionComputer.computeBytesForSigning(transaction);
        }, "The `chainID` field is not set");

        transaction.chainID = "localnet";
        transaction.version = 1;
        transaction.options = 2;

        assert.throws(() => {
            transactionComputer.computeBytesForSigning(transaction);
        }, `Non-empty transaction options requires transaction version >= ${MIN_TRANSACTION_VERSION_THAT_SUPPORTS_OPTIONS}`);

        transactionComputer.applyOptionsForHashSigning(transaction);

        assert.equal(transaction.version, 2);
        assert.equal(transaction.options, 3);
    });

    it("should compute bytes to verify transaction signature", async () => {
        let transaction = new Transaction({
            sender: alice.address,
            receiver: bob.address,
            gasLimit: 50000n,
            chainID: "D",
            nonce: 7n,
        });

        transaction.signature = await alice.signTransaction(transaction);

        const isSignedByAlice = await alice.verifyTransactionSignature(transaction, transaction.signature);

        const isSignedByBob = await bob.verifyTransactionSignature(transaction, transaction.signature);

        assert.equal(isSignedByAlice, true);
        assert.equal(isSignedByBob, false);
    });

    it("should compute bytes to verify transaction signature (signed by hash)", async () => {
        let transaction = new Transaction({
            sender: alice.address,
            receiver: bob.address,
            gasLimit: 50000n,
            chainID: "D",
            nonce: 7n,
        });

        transactionComputer.applyOptionsForHashSigning(transaction);

        transaction.signature = await alice.sign(transactionComputer.computeHashForSigning(transaction));

        const isSignedByAlice = await alice.verifyTransactionSignature(transaction, transaction.signature);

        const isSignedByBob = await bob.verifyTransactionSignature(transaction, transaction.signature);
        assert.equal(isSignedByAlice, true);
        assert.equal(isSignedByBob, false);
    });

    it("converts transaction to plain object and back", () => {
        const transaction = new Transaction({
            nonce: 90n,
            value: BigInt("123456789000000000000000000000"),
            sender: Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            receiver: Address.newFromBech32("drt1yxgvusarlvccy8p30gpl0gywwegd7nqghjqguyn672qwdfdg7y6qqta3dj"),
            senderUsername: "alice",
            receiverUsername: "bob",
            gasPrice: 1000000000n,
            gasLimit: 80000n,
            data: Buffer.from("hello"),
            chainID: "localnet",
            version: 2,
        });

        const plainObject = transaction.toPlainObject();
        const restoredTransaction = Transaction.newFromPlainObject(plainObject);

        assert.deepEqual(plainObject, transaction.toPlainObject());
        assert.deepEqual(restoredTransaction, Transaction.newFromPlainObject(plainObject));
        assert.deepEqual(restoredTransaction, transaction);
        assert.deepEqual(plainObject, {
            nonce: 90,
            value: "123456789000000000000000000000",
            sender: "drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh",
            receiver: "drt1yxgvusarlvccy8p30gpl0gywwegd7nqghjqguyn672qwdfdg7y6qqta3dj",
            senderUsername: "YWxpY2U=",
            receiverUsername: "Ym9i",
            gasPrice: 1000000000,
            gasLimit: 80000,
            data: "aGVsbG8=",
            chainID: "localnet",
            version: 2,
            options: undefined,
            guardian: undefined,
            relayer: undefined,
            signature: undefined,
            guardianSignature: undefined,
            relayerSignature: undefined,
        });
    });
    it("should serialize transaction with relayer", async () => {
        const transaction = new Transaction({
            chainID: networkConfig.chainID,
            sender: alice.address,
            receiver: alice.address,
            relayer: bob.address,
            gasLimit: 50000n,
            value: 0n,
            version: 2,
            nonce: 89n,
        });

        const serializedTransactionBytes = transactionComputer.computeBytesForSigning(transaction);
        const serializedTransaction = Buffer.from(serializedTransactionBytes).toString();

        assert.equal(
            serializedTransaction,
            `{"nonce":89,"value":"0","receiver":"drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh","sender":"drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh","gasPrice":1000000000,"gasLimit":50000,"chainID":"D","version":2,"relayer":"drt18h03w0y7qtqwtra3u4f0gu7e3kn2fslj83lqxny39m5c4rwaectswerhd2"}`,
        );
    });

    it("should test relayed v3", async () => {
        const transaction = new Transaction({
            chainID: networkConfig.chainID,
            sender: alice.address,
            receiver: alice.address,
            senderUsername: "alice",
            receiverUsername: "bob",
            gasLimit: 80000n,
            value: 0n,
            version: 2,
            nonce: 89n,
            data: Buffer.from("hello"),
        });

        assert.isFalse(transactionComputer.isRelayedV3Transaction(transaction));
        transaction.relayer = carol.address;
        assert.isTrue(transactionComputer.isRelayedV3Transaction(transaction));
    });
});
