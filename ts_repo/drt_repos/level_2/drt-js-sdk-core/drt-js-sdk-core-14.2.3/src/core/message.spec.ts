import { assert } from "chai";
import { Account } from "../accounts";
import { getTestWalletsPath } from "../testutils/utils";
import { DEFAULT_MESSAGE_VERSION, SDK_JS_SIGNER, UNKNOWN_SIGNER } from "./constants";
import { Message, MessageComputer } from "./message";

describe("test message", () => {
    let alice: Account;
    const messageComputer = new MessageComputer();

    before(async function () {
        alice = await Account.newFromPem(`${getTestWalletsPath()}/alice.pem`);
    });

    it("should test message compute bytes for signing", async () => {
        const data = Buffer.from("test message");

        const message = new Message({
            data: data,
        });

        const serialized = messageComputer.computeBytesForSigning(message);

        assert.equal(
            Buffer.from(serialized).toString("hex"),
            "0f6fce3fa6130fc58a25eaff6e157ea1bcb02fbf9773dca514dfaf3cd1e0bdfe",
        );
    });

    it("should create, sign, pack, unpack and verify message", async () => {
        const data = Buffer.from("test");

        const message = new Message({
            data: data,
            address: alice.address,
        });

        message.signature = await alice.signMessage(message);

        assert.equal(
            Buffer.from(message.signature).toString("hex"),
            "2758dcd8e5b5f7ccfaeaf7e0362ace8f83145c573da3f495d9218b2de367d161e007f3db5f288123b937cd7abaefaf7e66f417cc59ef8618c405de22578e6608",
        );

        const packedMessage = messageComputer.packMessage(message);
        assert.deepEqual(packedMessage, {
            address: "drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh",
            message: "74657374",
            signature:
                "2758dcd8e5b5f7ccfaeaf7e0362ace8f83145c573da3f495d9218b2de367d161e007f3db5f288123b937cd7abaefaf7e66f417cc59ef8618c405de22578e6608",
            version: 1,
            signer: SDK_JS_SIGNER,
        });

        const unpackedMessage = messageComputer.unpackMessage(packedMessage);
        assert.deepEqual(unpackedMessage.address, alice.address);
        assert.deepEqual(unpackedMessage.data, message.data);
        assert.deepEqual(unpackedMessage.signature, message.signature);
        assert.deepEqual(unpackedMessage.version, message.version);
        assert.deepEqual(unpackedMessage.signer, message.signer);

        const isValid = await alice.verifyMessageSignature(unpackedMessage, Buffer.from(unpackedMessage.signature!));
        assert.isTrue(isValid);
    });

    it("should unpack legacy message", async () => {
        const legacyMessage = {
            address: "drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh",
            message: "0x7468697320697320612074657374206d657373616765",
            signature:
                "0xb16847437049986f936dd4a0917c869730cbf29e40a0c0821ca70db33f44758c3d41bcbea446dee70dea13d50942343bb78e74979dc434bbb2b901e0f4fd1809",
            version: 1,
            signer: "DrtJS",
        };

        const message = messageComputer.unpackMessage(legacyMessage);
        assert.deepEqual(message.address, alice.address);
        assert.deepEqual(Buffer.from(message.data).toString(), "this is a test message");
        assert.deepEqual(
            Buffer.from(message.signature!).toString("hex"),
            "b16847437049986f936dd4a0917c869730cbf29e40a0c0821ca70db33f44758c3d41bcbea446dee70dea13d50942343bb78e74979dc434bbb2b901e0f4fd1809",
        );
        assert.deepEqual(message.version, DEFAULT_MESSAGE_VERSION);
        assert.equal(message.signer, "DrtJS");
    });

    it("should unpack message", async () => {
        const packedMessage = {
            address: "drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh",
            message: "0x7468697320697320612074657374206d657373616765",
            signature:
                "0xb16847437049986f936dd4a0917c869730cbf29e40a0c0821ca70db33f44758c3d41bcbea446dee70dea13d50942343bb78e74979dc434bbb2b901e0f4fd1809",
        };

        const message = messageComputer.unpackMessage(packedMessage);
        assert.deepEqual(message.address, alice.address);
        assert.deepEqual(Buffer.from(message.data).toString(), "this is a test message");
        assert.deepEqual(
            Buffer.from(message.signature!).toString("hex"),
            "b16847437049986f936dd4a0917c869730cbf29e40a0c0821ca70db33f44758c3d41bcbea446dee70dea13d50942343bb78e74979dc434bbb2b901e0f4fd1809",
        );
        assert.deepEqual(message.version, DEFAULT_MESSAGE_VERSION);
        assert.equal(message.signer, UNKNOWN_SIGNER);
    });
});
