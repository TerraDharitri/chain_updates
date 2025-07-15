import { assert } from "chai";
import { Account } from "../accounts";
import { Address, TokenTransfer, Transaction } from "../core";
import { getTestWalletsPath } from "../testutils";
import { ProtoSerializer } from "./serializer";

describe("serialize transactions", () => {
    let serializer = new ProtoSerializer();
    let alice: Account;
    let bob: Account;
    let carol: Account;
    before(async function () {
        alice = await Account.newFromPem(`${getTestWalletsPath()}/alice.pem`);
        bob = await Account.newFromPem(`${getTestWalletsPath()}/bob.pem`);
        carol = await Account.newFromPem(`${getTestWalletsPath()}/carol.pem`);
    });

    it("with no data, no value", async () => {
        let transaction = new Transaction({
            nonce: 89n,
            value: 0n,
            sender: alice.address,
            receiver: bob.address,
            gasLimit: 50000n,
            chainID: "local-testnet",
        });

        transaction.signature = await alice.signTransaction(transaction);

        let buffer = serializer.serializeTransaction(transaction);
        assert.equal(
            buffer.toString("hex"),
            "0859120200001a203ddf173c9e02c0e58fb1e552f473d98da6a4c3f23c7e034c912ee98a8dddce172a20391f932707a9dfa86d3bcbb3d5d0cc9f25ad0e680fe499f107d844b7e6ea71d5388094ebdc0340d08603520d6c6f63616c2d746573746e657458026240129faad6b4222905addd10a2b49941f0a0f559de52e78af98752916cd470980527f4822a8503ded124f6a54f2801c2617c1922a6d0d2c35c33daa437311a4903",
        );
    });

    it("with data, no value", async () => {
        let transaction = new Transaction({
            nonce: 90n,
            value: 0n,
            sender: alice.address,
            receiver: bob.address,
            gasLimit: 80000n,
            data: Buffer.from("hello"),
            chainID: "local-testnet",
        });

        transaction.signature = await alice.signTransaction(transaction);

        let buffer = serializer.serializeTransaction(transaction);
        assert.equal(
            buffer.toString("hex"),
            "085a120200001a203ddf173c9e02c0e58fb1e552f473d98da6a4c3f23c7e034c912ee98a8dddce172a20391f932707a9dfa86d3bcbb3d5d0cc9f25ad0e680fe499f107d844b7e6ea71d5388094ebdc034080f1044a0568656c6c6f520d6c6f63616c2d746573746e65745802624032290b79b021945a6650c5f38bcc8bf70f025323598396eaae68447912a84b5bc07897cf5261c48de1b0ab588148811cc394a9462c65372c8dfed1e245951807",
        );
    });

    it("with data, with value", async () => {
        let transaction = new Transaction({
            nonce: 91n,
            value: TokenTransfer.newFromNativeAmount(10000000000000000000n).amount,
            sender: alice.address,
            receiver: bob.address,
            gasLimit: 100000n,
            data: Buffer.from("for the book"),
            chainID: "local-testnet",
        });

        transaction.signature = await alice.signTransaction(transaction);

        let buffer = serializer.serializeTransaction(transaction);
        assert.equal(
            buffer.toString("hex"),
            "085b1209008ac7230489e800001a203ddf173c9e02c0e58fb1e552f473d98da6a4c3f23c7e034c912ee98a8dddce172a20391f932707a9dfa86d3bcbb3d5d0cc9f25ad0e680fe499f107d844b7e6ea71d5388094ebdc0340a08d064a0c666f722074686520626f6f6b520d6c6f63616c2d746573746e6574580262402605d0c7298a218608fe5324cb9bdea34f52b301b73dc8a351c902220e880c50ad2b4039e68ac403103b24eaa4e8bc32fbfa90a58da6dce818b8491929f6550d",
        );
    });

    it("with data, with large value", async () => {
        let transaction = new Transaction({
            nonce: 92n,
            value: 123456789000000000000000000000n,
            sender: alice.address,
            receiver: bob.address,
            gasLimit: 100000n,
            data: Buffer.from("for the spaceship"),
            chainID: "local-testnet",
        });

        transaction.signature = await alice.signTransaction(transaction);

        let buffer = serializer.serializeTransaction(transaction);
        assert.equal(
            buffer.toString("hex"),
            "085c120e00018ee90ff6181f3761632000001a203ddf173c9e02c0e58fb1e552f473d98da6a4c3f23c7e034c912ee98a8dddce172a20391f932707a9dfa86d3bcbb3d5d0cc9f25ad0e680fe499f107d844b7e6ea71d5388094ebdc0340a08d064a11666f722074686520737061636573686970520d6c6f63616c2d746573746e65745802624061deaf3566fc651e5d50ae8f8989771c149077b8f721e02e886be66f013be4e4ae1eadc6d36d54ce0d3899b34db6e6198d01f7b155243305fb49199f21937609",
        );
    });

    it("with nonce = 0", async () => {
        let transaction = new Transaction({
            nonce: 0n,
            value: 0n,
            sender: alice.address,
            receiver: bob.address,
            gasLimit: 80000n,
            data: Buffer.from("hello"),
            chainID: "local-testnet",
            version: 1,
        });

        transaction.signature = await alice.signTransaction(transaction);

        let buffer = serializer.serializeTransaction(transaction);
        assert.equal(
            buffer.toString("hex"),
            "120200001a203ddf173c9e02c0e58fb1e552f473d98da6a4c3f23c7e034c912ee98a8dddce172a20391f932707a9dfa86d3bcbb3d5d0cc9f25ad0e680fe499f107d844b7e6ea71d5388094ebdc034080f1044a0568656c6c6f520d6c6f63616c2d746573746e65745801624024b30a91becdbba0dddc74886d52253a221ed2c381da0c6d7ea4c6482075b577c1b105c588ba8ec9442182b345e3cda975ecfc3785481ad44eb2dfd0464aee00",
        );
    });

    it("with usernames", async () => {
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

        const buffer = serializer.serializeTransaction(transaction);
        assert.equal(
            buffer.toString("hex"),
            "08cc011209000de0b6b3a76400001a20391f932707a9dfa86d3bcbb3d5d0cc9f25ad0e680fe499f107d844b7e6ea71d52205616c6963652a20b2a11555ce521e4944e09ab17549d85b487dcd26c84b5017a39e31a3670889ba32056361726f6c388094ebdc0340d0860352015458026240c6aee90d57933382f0b06ce75619b7a4dfaa05f22f885dee7272654ea6bf96b66f264ab8fbe84282a0817893c046e9f2dd875621655b3215b4fc6ba251ccd808",
        );
    });
});
