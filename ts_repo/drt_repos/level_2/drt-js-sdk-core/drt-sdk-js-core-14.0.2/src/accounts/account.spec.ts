import { assert } from "chai";
import { Address, Message, Transaction } from "../core";
import { getTestWalletsPath } from "../testutils/utils";
import { KeyPair, UserSecretKey } from "../wallet";
import { Account } from "./account";

describe("test account methods", function () {
    const DUMMY_MNEMONIC =
        "bread type ride autumn corn maid benefit pole that normal orchard confirm napkin degree arrow guitar offer you enjoy bronze more onion push decorate";
    const alice = `${getTestWalletsPath()}/alice.pem`;
    it("should create account from pem file", async function () {
        const account = await Account.newFromPem(alice);

        assert.equal(
            account.secretKey.valueOf().toString("hex"),
            "2bbcdae7e193924fa0d301e7a12c7defc92a93bc5e587cc968f04fcb86022e1c",
        );
        assert.equal(account.address.toBech32(), "drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
    });

    it("should create account from keystore", async function () {
        const account = Account.newFromKeystore(`${getTestWalletsPath()}/withDummyMnemonic.json`, "password");

        assert.equal(
            account.secretKey.valueOf().toString("hex"),
            "2bbcdae7e193924fa0d301e7a12c7defc92a93bc5e587cc968f04fcb86022e1c",
        );
        assert.equal(account.address.toBech32(), "drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
    });

    it("should create account from mnemonic", async function () {
        const account = Account.newFromMnemonic(DUMMY_MNEMONIC);

        assert.equal(
            account.secretKey.valueOf().toString("hex"),
            "2bbcdae7e193924fa0d301e7a12c7defc92a93bc5e587cc968f04fcb86022e1c",
        );
        assert.equal(account.address.toBech32(), "drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
    });

    it("should create account from keypair", async function () {
        const secretKey = UserSecretKey.fromString("391f932707a9dfa86d3bcbb3d5d0cc9f25ad0e680fe499f107d844b7e6ea71d5");
        const keypair = new KeyPair(secretKey);
        const account = Account.newFromKeypair(keypair);
        assert.deepEqual(account.secretKey, secretKey);
        assert.equal(account.address.toBech32(), "drt1vk5umh55umf3qajlv8pkvraramgeafq0fnav5flusjq0xafrypnqsjm9qh");
    });

    it("should increase nonce on account", async function () {
        const account = Account.newFromMnemonic(DUMMY_MNEMONIC);
        account.nonce = 42n;

        assert.equal(account.getNonceThenIncrement(), 42n);
        assert.equal(account.getNonceThenIncrement(), 43n);
    });

    it("should sign transaction", async function () {
        const transaction = new Transaction({
            nonce: 89n,
            value: 0n,
            receiver: Address.newFromBech32("drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c"),
            sender: Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            gasPrice: 1000000000n,
            gasLimit: 50000n,
            data: new Uint8Array(),
            chainID: "local-testnet",
            version: 1,
            options: 0,
        });

        const account = Account.newFromMnemonic(DUMMY_MNEMONIC);
        transaction.signature = await account.signTransaction(transaction);

        assert.equal(
            Buffer.from(transaction.signature).toString("hex"),
            "d43261cf310bd2c1e5d7cc4eebca4dfbb8afd04e7645179ee31cf879f88b9dc23ff64919fbba0693229fa133973ef324904104de41f05801100e2df9e8bbb904",
        );
    });

    it("should sign message", async function () {
        const message = new Message({
            data: new Uint8Array(Buffer.from("hello")),
        });

        const account = Account.newFromMnemonic(DUMMY_MNEMONIC);
        message.signature = await account.signMessage(message);

        assert.equal(
            Buffer.from(message.signature).toString("hex"),
            "33edba0c691b5a3e8211a5fa63508a4f0c5ba7ac066413ea660e8ec9145c57521d13c304f21bd4687e9f4e118c8df0df6d20ad59b56dffbf749dd2b3b377740f",
        );
    });

    it("should verify message", async function () {
        const message = new Message({
            data: new Uint8Array(Buffer.from("hello")),
        });

        const account = Account.newFromMnemonic(DUMMY_MNEMONIC);
        message.signature = await account.signMessage(message);
        const isVerified = await account.verifyMessageSignature(message, message.signature);

        assert.isTrue(isVerified);
    });

    it("should sign and verify transaction", async function () {
        const transaction = new Transaction({
            nonce: 89n,
            value: 0n,
            receiver: Address.newFromBech32("drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c"),
            sender: Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            gasPrice: 1000000000n,
            gasLimit: 50000n,
            data: new Uint8Array(),
            chainID: "local-testnet",
            version: 1,
            options: 0,
        });

        const account = Account.newFromMnemonic(DUMMY_MNEMONIC);
        transaction.signature = await account.signTransaction(transaction);

        assert.equal(
            Buffer.from(transaction.signature).toString("hex"),
            "d43261cf310bd2c1e5d7cc4eebca4dfbb8afd04e7645179ee31cf879f88b9dc23ff64919fbba0693229fa133973ef324904104de41f05801100e2df9e8bbb904",
        );

        const isVerified = await account.verifyTransactionSignature(transaction, transaction.signature);
        assert.isTrue(isVerified);
    });
});
