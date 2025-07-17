import { assert } from "chai";
import { Address, Transaction, TransactionComputer } from "../core";
import { KeyPair } from "./keypair";
import { UserSecretKey } from "./userKeys";

describe("test keypair", () => {
    it("should create keypair", () => {
        const buffer_hex = "2bbcdae7e193924fa0d301e7a12c7defc92a93bc5e587cc968f04fcb86022e1c";
        const buffer = Uint8Array.from(Buffer.from(buffer_hex, "hex"));
        const userSecretKey = UserSecretKey.fromString(buffer_hex);
        let keypair = KeyPair.newFromBytes(buffer);
        let secretKey = keypair.getSecretKey();
        assert.equal(secretKey.hex(), buffer_hex);
        assert.equal(keypair.secretKey.hex(), buffer_hex);
        assert.deepEqual(secretKey, userSecretKey);

        keypair = new KeyPair(secretKey);
        assert.deepEqual(keypair.getSecretKey(), userSecretKey);
        assert.deepEqual(keypair.getPublicKey(), userSecretKey.generatePublicKey());

        keypair = KeyPair.generate();
        const pubkey = keypair.getPublicKey();
        secretKey = keypair.getSecretKey();
        assert.lengthOf(pubkey.valueOf(), 32);
        assert.lengthOf(secretKey.valueOf(), 32);
    });

    it("should sign and verify transaction", async () => {
        const transaction = new Transaction({
            nonce: 89n,
            value: 0n,
            receiver: Address.newFromBech32("drt1yxgvusarlvccy8p30gpl0gywwegd7nqghjqguyn672qwdfdg7y6qqta3dj"),
            sender: Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            gasPrice: 1000000000n,
            gasLimit: 50000n,
            chainID: "local-testnet",
            version: 1,
            options: 0,
        });
        const bufferHex = "2bbcdae7e193924fa0d301e7a12c7defc92a93bc5e587cc968f04fcb86022e1c";
        const buffer = Uint8Array.from(Buffer.from(bufferHex, "hex"));
        const keypair = KeyPair.newFromBytes(buffer);

        const transactionComputer = new TransactionComputer();
        const serializedTx = transactionComputer.computeBytesForSigning(transaction);
        transaction.signature = await keypair.sign(serializedTx);
        assert.equal(
            Buffer.from(transaction.signature).toString("hex"),
            "cd69f88f09cdeb934fcdef5c711fd2d6ecc080c66208b04c1d28d699e101cbed2eab0d942bd47157fc772603bc639cc152a278b8f89a45ab50df4f8e79a1d90b",
        );
        assert.isTrue(await keypair.verify(serializedTx, transaction.signature));
    });
});
