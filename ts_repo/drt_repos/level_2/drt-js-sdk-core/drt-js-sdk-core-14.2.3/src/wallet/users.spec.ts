import { assert } from "chai";
import path from "path";
import { Account } from "../accounts";
import { Address, ErrBadMnemonicEntropy, ErrInvariantFailed, Message, Transaction } from "../core";
import {
    DummyMnemonicOf12Words,
    loadMnemonic,
    loadPassword,
    loadTestKeystore,
    loadTestWallet,
    TestWallet,
} from "./../testutils/wallets";
import { Randomness } from "./crypto";
import { Mnemonic } from "./mnemonic";
import { UserSecretKey } from "./userKeys";
import { UserSigner } from "./userSigner";
import { UserVerifier } from "./userVerifier";
import { UserWallet } from "./userWallet";

describe("test user wallets", () => {
    let alice: TestWallet, bob: TestWallet, carol: TestWallet;
    let password: string;
    let dummyMnemonic: string;

    before(async function () {
        alice = await loadTestWallet("alice");
        bob = await loadTestWallet("bob");
        carol = await loadTestWallet("carol");
        password = await loadPassword();
        dummyMnemonic = await loadMnemonic();
    });

    it("should generate mnemonic", () => {
        let mnemonic = Mnemonic.generate();
        let words = mnemonic.getWords();
        assert.lengthOf(words, 24);
    });

    it("should convert entropy to mnemonic and back", () => {
        function testConversion(text: string, entropyHex: string) {
            const entropyFromMnemonic = Mnemonic.fromString(text).getEntropy();
            const mnemonicFromEntropy = Mnemonic.fromEntropy(Buffer.from(entropyHex, "hex"));

            assert.equal(Buffer.from(entropyFromMnemonic).toString("hex"), entropyHex);
            assert.equal(mnemonicFromEntropy.toString(), text);
        }

        testConversion(
            "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
            "00000000000000000000000000000000",
        );

        testConversion(
            "bread type ride autumn corn maid benefit pole that normal orchard confirm napkin degree arrow guitar offer you enjoy bronze more onion push decorate",
            "1b7d7ae587d3070c054d3cdff2c67017792e73832b3e997fe52a8e58fd35aba9",
        );

        assert.throws(
            () => {
                Mnemonic.fromEntropy(Buffer.from("abba", "hex"));
            },
            ErrBadMnemonicEntropy,
            `Bad mnemonic entropy`,
        );
    });

    it("should derive keys", async () => {
        let mnemonic = Mnemonic.fromString(dummyMnemonic);

        assert.equal(mnemonic.deriveKey(0).hex(), alice.secretKeyHex);
        assert.equal(mnemonic.deriveKey(1).hex(), "1f4d9984ff57a9bcc7b8aea32069e41d36366e4dd9e08f55c6691168de06f2c3");
        assert.equal(mnemonic.deriveKey(2).hex(), "6c030765ecd8dce0e8aa8e15ab10823d5ae5dc682d3cb6c260640f01def7a587");
    });

    it("should derive keys (12 words)", async () => {
        const mnemonic = Mnemonic.fromString(DummyMnemonicOf12Words);

        assert.equal(
            mnemonic.deriveKey(0).generatePublicKey().toAddress().toBech32(),
            "drt1l8g9dk3gz035gkjhwegsjkqzdu3augrwhcfxrnucnyyrpc2220pq4flasr",
        );
        assert.equal(
            mnemonic.deriveKey(1).generatePublicKey().toAddress().toBech32(),
            "drt1fmhwg84rldg0xzngf53m0y607wvefvamh07n2mkypedx27lcqntsg78vxl",
        );
        assert.equal(
            mnemonic.deriveKey(2).generatePublicKey().toAddress().toBech32(),
            "drt1tyuyemt4xz2yjvc7rxxp8kyfmk2n3h8gv3aavzd9ru4v2vhrkcksuhwdgv",
        );

        assert.equal(
            mnemonic.deriveKey(0).generatePublicKey().toAddress("test").toBech32(),
            "test1l8g9dk3gz035gkjhwegsjkqzdu3augrwhcfxrnucnyyrpc2220pqc6tnnf",
        );
        assert.equal(
            mnemonic.deriveKey(1).generatePublicKey().toAddress("xdrt").toBech32(),
            "xdrt1fmhwg84rldg0xzngf53m0y607wvefvamh07n2mkypedx27lcqnts0f2w3t",
        );
        assert.equal(
            mnemonic.deriveKey(2).generatePublicKey().toAddress("ydrt").toBech32(),
            "ydrt1tyuyemt4xz2yjvc7rxxp8kyfmk2n3h8gv3aavzd9ru4v2vhrkckswmkvs2",
        );
    });

    it("should create secret key", () => {
        const keyHex = alice.secretKeyHex;
        const fromBuffer = new UserSecretKey(Buffer.from(keyHex, "hex"));
        const fromArray = new UserSecretKey(Uint8Array.from(Buffer.from(keyHex, "hex")));
        const fromHex = UserSecretKey.fromString(keyHex);

        assert.equal(fromBuffer.hex(), keyHex);
        assert.equal(fromArray.hex(), keyHex);
        assert.equal(fromHex.hex(), keyHex);
    });

    it("should compute public key (and address)", () => {
        let secretKey: UserSecretKey;

        secretKey = new UserSecretKey(Buffer.from(alice.secretKeyHex, "hex"));
        assert.equal(secretKey.generatePublicKey().hex(), alice.address.hex());
        assert.deepEqual(secretKey.generatePublicKey().toAddress(), alice.address);

        secretKey = new UserSecretKey(Buffer.from(bob.secretKeyHex, "hex"));
        assert.equal(secretKey.generatePublicKey().hex(), bob.address.hex());
        assert.deepEqual(secretKey.generatePublicKey().toAddress(), bob.address);

        secretKey = new UserSecretKey(Buffer.from(carol.secretKeyHex, "hex"));
        assert.equal(secretKey.generatePublicKey().hex(), carol.address.hex());
        assert.deepEqual(secretKey.generatePublicKey().toAddress(), carol.address);
    });

    it("should throw error when invalid input", () => {
        assert.throw(() => new UserSecretKey(Buffer.alloc(42)), ErrInvariantFailed);
        assert.throw(() => UserSecretKey.fromString("foobar"), ErrInvariantFailed);
    });

    it("should handle PEM files", () => {
        assert.equal(UserSecretKey.fromPem(alice.pemFileText).hex(), alice.secretKeyHex);
        assert.equal(UserSecretKey.fromPem(bob.pemFileText).hex(), bob.secretKeyHex);
        assert.equal(UserSecretKey.fromPem(carol.pemFileText).hex(), carol.secretKeyHex);
    });

    it("should create and load keystore files (with secret keys)", function () {
        this.timeout(10000);

        let aliceSecretKey = UserSecretKey.fromString(alice.secretKeyHex);
        let bobSecretKey = UserSecretKey.fromString(bob.secretKeyHex);
        let carolSecretKey = UserSecretKey.fromString(carol.secretKeyHex);

        console.time("encrypt");
        let aliceKeyFile = UserWallet.fromSecretKey({ secretKey: aliceSecretKey, password: password });
        let bobKeyFile = UserWallet.fromSecretKey({ secretKey: bobSecretKey, password: password });
        let carolKeyFile = UserWallet.fromSecretKey({ secretKey: carolSecretKey, password: password });
        console.timeEnd("encrypt");

        assert.equal(aliceKeyFile.toJSON().bech32, alice.address.toBech32());
        assert.equal(bobKeyFile.toJSON().bech32, bob.address.toBech32());
        assert.equal(carolKeyFile.toJSON().bech32, carol.address.toBech32());

        console.time("decrypt");
        assert.deepEqual(UserWallet.decryptSecretKey(aliceKeyFile.toJSON(), password), aliceSecretKey);
        assert.deepEqual(UserWallet.decryptSecretKey(bobKeyFile.toJSON(), password), bobSecretKey);
        assert.deepEqual(UserWallet.decryptSecretKey(carolKeyFile.toJSON(), password), carolSecretKey);
        console.timeEnd("decrypt");

        // With provided randomness, in order to reproduce our development wallets

        aliceKeyFile = UserWallet.fromSecretKey({
            secretKey: aliceSecretKey,
            password: password,
            randomness: new Randomness({
                id: alice.keyFileObject.id,
                iv: Buffer.from(alice.keyFileObject.crypto.cipherparams.iv, "hex"),
                salt: Buffer.from(alice.keyFileObject.crypto.kdfparams.salt, "hex"),
            }),
        });

        bobKeyFile = UserWallet.fromSecretKey({
            secretKey: bobSecretKey,
            password: password,
            randomness: new Randomness({
                id: bob.keyFileObject.id,
                iv: Buffer.from(bob.keyFileObject.crypto.cipherparams.iv, "hex"),
                salt: Buffer.from(bob.keyFileObject.crypto.kdfparams.salt, "hex"),
            }),
        });

        carolKeyFile = UserWallet.fromSecretKey({
            secretKey: carolSecretKey,
            password: password,
            randomness: new Randomness({
                id: carol.keyFileObject.id,
                iv: Buffer.from(carol.keyFileObject.crypto.cipherparams.iv, "hex"),
                salt: Buffer.from(carol.keyFileObject.crypto.kdfparams.salt, "hex"),
            }),
        });

        assert.deepEqual(aliceKeyFile.toJSON(), alice.keyFileObject);
        assert.deepEqual(bobKeyFile.toJSON(), bob.keyFileObject);
        assert.deepEqual(carolKeyFile.toJSON(), carol.keyFileObject);
    });

    it("should load keystore files (with secret keys, but without 'kind' field)", async function () {
        const keyFileObject = await loadTestKeystore("withoutKind.json");
        const secretKey = UserWallet.decryptSecretKey(keyFileObject, password);

        assert.equal(
            secretKey.generatePublicKey().toAddress().toBech32(),
            "drt18etzjgpc7h9xr5hgm62qsyvrvunl43htr5p4dp29c79pl4nk0essfjlk83",
        );
    });

    it("should create and load keystore files (with mnemonics)", async function () {
        this.timeout(10000);

        const wallet = UserWallet.fromMnemonic({ mnemonic: dummyMnemonic, password: password });
        const json = wallet.toJSON();

        assert.equal(json.version, 4);
        assert.equal(json.kind, "mnemonic");
        assert.isUndefined(json.toBech32);

        const mnemonic = UserWallet.decryptMnemonic(json, password);
        const mnemonicText = mnemonic.toString();

        assert.equal(mnemonicText, dummyMnemonic);
        assert.equal(mnemonic.deriveKey(0).generatePublicKey().toAddress().toBech32(), alice.address.toBech32());
        assert.equal(mnemonic.deriveKey(1).generatePublicKey().toAddress().toBech32(), "drt1tzkwpg0et0s7fp46a7je9h0gv2v55t9mamqrhgja7wypp4yf5d0se3nzyj");
        assert.equal(mnemonic.deriveKey(2).generatePublicKey().toAddress().toBech32(), "drt15hpu70r43r3hx9evqmu2z04097ye5t0jgrw3lhxw5rnge0k89nlsh83ydx");

        // With provided randomness, in order to reproduce our test wallets
        const expectedDummyWallet = await loadTestKeystore("withDummyMnemonic.json");
        const dummyWallet = UserWallet.fromMnemonic({
            mnemonic: dummyMnemonic,
            password: password,
            randomness: new Randomness({
                id: "22efac61-f898-4682-9b96-2e4fb49f56d5",
                iv: Buffer.from("69f26c6e4181ebfbce9de4080122af15", "hex"),
                salt: Buffer.from("a3c74d0544a697f3d032da820627f95db417f13a86ec8468bba2acd116c34adf", "hex"),
            }),
        });

        assert.deepEqual(dummyWallet.toJSON(), expectedDummyWallet);
    });

    it("should create user wallet from secret key, but without 'kind' field", async function () {
        const keyFileObject = await loadTestKeystore("withoutKind.json");
        const secretKey = UserWallet.decrypt(keyFileObject, password);

        assert.equal(
            secretKey.generatePublicKey().toAddress().toBech32(),
            "drt18etzjgpc7h9xr5hgm62qsyvrvunl43htr5p4dp29c79pl4nk0essfjlk83",
        );
    });

    it("should loadSecretKey, but without 'kind' field", async function () {
        const testdataPath = path.resolve(__dirname, "..", "testdata/testwallets");
        const keystorePath = path.resolve(testdataPath, "withoutKind.json");
        const secretKey = UserWallet.loadSecretKey(keystorePath, password);

        assert.equal(
            secretKey.generatePublicKey().toAddress().toBech32(),
            "drt18etzjgpc7h9xr5hgm62qsyvrvunl43htr5p4dp29c79pl4nk0essfjlk83",
        );
    });

    it("should throw when calling loadSecretKey with unecessary address index", async function () {
        const keyFileObject = await loadTestKeystore("alice.json");

        assert.throws(
            () => UserWallet.decrypt(keyFileObject, password, 42),
            "addressIndex must not be provided when kind == 'secretKey'",
        );
    });

    it("should loadSecretKey with mnemonic", async function () {
        const keyFileObject = await loadTestKeystore("withDummyMnemonic.json");

        assert.equal(
            UserWallet.decrypt(keyFileObject, password, 0).generatePublicKey().toAddress().toBech32(),
            "drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh",
        );
        assert.equal(
            UserWallet.decrypt(keyFileObject, password, 1).generatePublicKey().toAddress().toBech32(),
            "drt1tzkwpg0et0s7fp46a7je9h0gv2v55t9mamqrhgja7wypp4yf5d0se3nzyj",
        );
        assert.equal(
            UserWallet.decrypt(keyFileObject, password, 2).generatePublicKey().toAddress().toBech32(),
            "drt15hpu70r43r3hx9evqmu2z04097ye5t0jgrw3lhxw5rnge0k89nlsh83ydx",
        );
    });

    it("should sign transactions", async () => {
        let signer = new Account(
            UserSecretKey.fromString("1a927e2af5306a9bb2ea777f73e06ecc0ac9aaa72fb4ea3fecf659451394cccf"),
        );
        let verifier = new UserVerifier(
            UserSecretKey.fromString(
                "1a927e2af5306a9bb2ea777f73e06ecc0ac9aaa72fb4ea3fecf659451394cccf",
            ).generatePublicKey(),
        );

        // With data field
        let transaction = new Transaction({
            nonce: 0n,
            value: 0n,
            sender: Address.newFromBech32("drt1l453hd0gt5gzdp7czpuall8ggt2dcv5zwmfdf3sd3lguxseux2fsxvluwu"),
            receiver: Address.newFromBech32("drt1cux02zersde0l7hhklzhywcxk4u9n4py5tdxyx7vrvhnza2r4gmqgsejha"),
            gasPrice: 1000000000n,
            gasLimit: 50000n,
            data: new TextEncoder().encode("foo"),
            chainID: "1",
        });

        let serialized = transaction.serializeForSigning();
        let signature = await signer.sign(serialized);

        assert.deepEqual(await signer.sign(serialized), await signer.sign(Uint8Array.from(serialized)));
        assert.deepEqual(
            serialized.toString(),
            `{"nonce":0,"value":"0","receiver":"drt1cux02zersde0l7hhklzhywcxk4u9n4py5tdxyx7vrvhnza2r4gmqgsejha","sender":"drt1l453hd0gt5gzdp7czpuall8ggt2dcv5zwmfdf3sd3lguxseux2fsxvluwu","gasPrice":1000000000,"gasLimit":50000,"data":"Zm9v","chainID":"1","version":2}`,
        );
        assert.equal(
            Buffer.from(signature).toString("hex"),
            "fba90410603f6a3d89f0faaee745eb97dc09de9a21bd020cd05687893bac4e800e01e8e32da31b15fdca483d422b03fda71e3285903313af56b35714a796ba01",
        );
        assert.isTrue(await verifier.verify(serialized, signature));

        // Without data field
        transaction = new Transaction({
            nonce: 8n,
            value: 10000000000000000000n,
            sender: Address.newFromBech32("drt1l453hd0gt5gzdp7czpuall8ggt2dcv5zwmfdf3sd3lguxseux2fsxvluwu"),
            receiver: Address.newFromBech32("drt1cux02zersde0l7hhklzhywcxk4u9n4py5tdxyx7vrvhnza2r4gmqgsejha"),
            gasPrice: 1000000000n,
            gasLimit: 50000n,
            chainID: "1",
        });

        serialized = transaction.serializeForSigning();
        signature = await signer.sign(serialized);

        assert.deepEqual(await signer.sign(serialized), await signer.sign(Uint8Array.from(serialized)));
        assert.equal(
            serialized.toString(),
            `{"nonce":8,"value":"10000000000000000000","receiver":"drt1cux02zersde0l7hhklzhywcxk4u9n4py5tdxyx7vrvhnza2r4gmqgsejha","sender":"drt1l453hd0gt5gzdp7czpuall8ggt2dcv5zwmfdf3sd3lguxseux2fsxvluwu","gasPrice":1000000000,"gasLimit":50000,"chainID":"1","version":2}`,
        );
        assert.equal(
            Buffer.from(signature).toString("hex"),
            "37ecf2f4ddf853e5bcd7c134f86894f88df0bfa4585c2c58017aa0d9d8eda5df45e9845cda8c6aabe1ba4f5d04603ccd89c688b56fae967cd24bca31df387002",
        );
    });

    it("guardian should sign transactions from PEM", async () => {
        // bob is the guardian
        let signer = new UserSigner(
            UserSecretKey.fromString("1a927e2af5306a9bb2ea777f73e06ecc0ac9aaa72fb4ea3fecf659451394cccf"),
        );
        let verifier = new UserVerifier(
            UserSecretKey.fromString(
                "1a927e2af5306a9bb2ea777f73e06ecc0ac9aaa72fb4ea3fecf659451394cccf",
            ).generatePublicKey(),
        );
        let guardianSigner = new UserSigner(UserSecretKey.fromPem(bob.pemFileText));

        // With data field
        let transaction = new Transaction({
            nonce: 0n,
            value: 0n,
            receiver: Address.newFromBech32("drt1cux02zersde0l7hhklzhywcxk4u9n4py5tdxyx7vrvhnza2r4gmqgsejha"),
            sender: Address.newFromBech32("drt1l453hd0gt5gzdp7czpuall8ggt2dcv5zwmfdf3sd3lguxseux2fsxvluwu"),
            gasPrice: 1000000000n,
            gasLimit: 50000n,
            data: new TextEncoder().encode("foo"),
            chainID: "1",
            guardian: Address.newFromBech32("drt1yxgvusarlvccy8p30gpl0gywwegd7nqghjqguyn672qwdfdg7y6qqta3dj"),
            options: 2,
            version: 2,
        });

        let serialized = transaction.serializeForSigning();
        let signature = await signer.sign(serialized);
        let guardianSignature = await guardianSigner.sign(serialized);

        assert.deepEqual(
            serialized.toString(),
            `{"nonce":0,"value":"0","receiver":"drt1cux02zersde0l7hhklzhywcxk4u9n4py5tdxyx7vrvhnza2r4gmqgsejha","sender":"drt1l453hd0gt5gzdp7czpuall8ggt2dcv5zwmfdf3sd3lguxseux2fsxvluwu","gasPrice":1000000000,"gasLimit":50000,"data":"Zm9v","chainID":"1","version":2,"options":2,"guardian":"drt1yxgvusarlvccy8p30gpl0gywwegd7nqghjqguyn672qwdfdg7y6qqta3dj"}`,
        );
        assert.equal(
            Buffer.from(signature).toString("hex"),
            "c690e12720f2cb1f61c8569b9f27564936a80bf4d6fb74147565b0cb8f5c98b083beb7216965fc2568e7e42cca65fcae214f33050b1a9a7bb05702c8040d6e05",
        );
        assert.equal(
            Buffer.from(guardianSignature).toString("hex"),
            "04a845a1af7f8f87e42d0162c4bbd0c4d74bdc7458dad2539aed67eaa82d619ce94706b8d7745224a8760d4bea2b328cfae57af30962246c5283c8cabe1aef00",
        );
        assert.isTrue(await verifier.verify(serialized, signature));

        // Without data field
        transaction = new Transaction({
            nonce: 8n,
            value: 10000000000000000000n,
            receiver: Address.newFromBech32("drt1cux02zersde0l7hhklzhywcxk4u9n4py5tdxyx7vrvhnza2r4gmqgsejha"),
            sender: Address.newFromBech32("drt1l453hd0gt5gzdp7czpuall8ggt2dcv5zwmfdf3sd3lguxseux2fsxvluwu"),
            gasPrice: 1000000000n,
            gasLimit: 50000n,
            chainID: "1",
            guardian: Address.newFromBech32("drt1yxgvusarlvccy8p30gpl0gywwegd7nqghjqguyn672qwdfdg7y6qqta3dj"),
            options: 2,
            version: 2,
        });

        serialized = transaction.serializeForSigning();
        signature = await signer.sign(serialized);
        guardianSignature = await guardianSigner.sign(serialized);

        assert.equal(
            serialized.toString(),
            `{"nonce":8,"value":"10000000000000000000","receiver":"drt1cux02zersde0l7hhklzhywcxk4u9n4py5tdxyx7vrvhnza2r4gmqgsejha","sender":"drt1l453hd0gt5gzdp7czpuall8ggt2dcv5zwmfdf3sd3lguxseux2fsxvluwu","gasPrice":1000000000,"gasLimit":50000,"chainID":"1","version":2,"options":2,"guardian":"drt1yxgvusarlvccy8p30gpl0gywwegd7nqghjqguyn672qwdfdg7y6qqta3dj"}`,
        );
        assert.equal(
            Buffer.from(signature).toString("hex"),
            "9dd1fac96ec98f085ec035d643357b8495dccd300bcb406281017b94e7c6a796341ef8470a7633536884f032e4966ba96396da403c1fc9bdd190cc2f87defd02",
        );
        assert.equal(
            Buffer.from(guardianSignature).toString("hex"),
            "084bbff3e27eb6c2be55ffb7b105df7b011bc732dd0bdb3163fe01050c0e0de276f7f19b4a7f7afb1743278d2142de306d54b56c5e277fabc493038f64f6fd0a",
        );
        assert.isTrue(await verifier.verify(serialized, signature));
    });

    it("should sign transactions using PEM files", async () => {
        const signer = UserSigner.fromPem(alice.pemFileText);

        const transaction = new Transaction({
            nonce: 0n,
            value: 0n,
            sender: signer.getAddress(),
            receiver: Address.newFromBech32("drt1cux02zersde0l7hhklzhywcxk4u9n4py5tdxyx7vrvhnza2r4gmqgsejha"),
            gasPrice: 1000000000n,
            gasLimit: 50000n,
            data: new TextEncoder().encode("foo"),
            chainID: "1",
        });

        const serialized = transaction.serializeForSigning();
        const signature = await signer.sign(serialized);

        assert.deepEqual(await signer.sign(serialized), await signer.sign(Uint8Array.from(serialized)));
        assert.equal(
            Buffer.from(signature).toString("hex"),
            "32eb4a8b969fc04083432628282859617abf2c670c95963f42346e9ff39aa8d2388a01508d37eded4b5793213e739398f340208fc3d4dddecb09e79035097b0b",
        );
    });

    it("signs a general message", async function () {
        let signer = new UserSigner(
            UserSecretKey.fromString("1a927e2af5306a9bb2ea777f73e06ecc0ac9aaa72fb4ea3fecf659451394cccf"),
        );
        let verifier = new UserVerifier(
            UserSecretKey.fromString(
                "1a927e2af5306a9bb2ea777f73e06ecc0ac9aaa72fb4ea3fecf659451394cccf",
            ).generatePublicKey(),
        );

        const message = new Message({
            data: new TextEncoder().encode(JSON.stringify({ foo: "hello", bar: "world" })),
        });

        const signature = await signer.sign(message.data);

        assert.deepEqual(await signer.sign(message.data), await signer.sign(Uint8Array.from(message.data)));
        assert.isTrue(await verifier.verify(message.data, signature));
        assert.isTrue(await verifier.verify(Uint8Array.from(message.data), Uint8Array.from(signature)));
        assert.isFalse(await verifier.verify(Buffer.from("hello"), signature));
        assert.isFalse(await verifier.verify(new TextEncoder().encode("hello"), signature));
    });

    it("should create UserSigner from wallet", async function () {
        const keyFileObjectWithoutKind = await loadTestKeystore("withoutKind.json");
        const keyFileObjectWithMnemonic = await loadTestKeystore("withDummyMnemonic.json");
        const keyFileObjectWithSecretKey = await loadTestKeystore("withDummySecretKey.json");

        assert.equal(
            UserSigner.fromWallet(keyFileObjectWithoutKind, password).getAddress().toBech32(),
            "drt18etzjgpc7h9xr5hgm62qsyvrvunl43htr5p4dp29c79pl4nk0essfjlk83",
        );
        assert.equal(
            UserSigner.fromWallet(keyFileObjectWithMnemonic, password).getAddress().toBech32(),
            "drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh",
        );
        assert.equal(
            UserSigner.fromWallet(keyFileObjectWithSecretKey, password).getAddress().toBech32(),
            "drt1pvhrxtm26zgc59lnky049dyj4s8q3snaeg26kmxegx90qgvzrczq92zp9q",
        );
        assert.equal(
            UserSigner.fromWallet(keyFileObjectWithMnemonic, password, 0).getAddress().toBech32(),
            "drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh",
        );
        assert.equal(
            UserSigner.fromWallet(keyFileObjectWithMnemonic, password, 1).getAddress().toBech32(),
            "drt1tzkwpg0et0s7fp46a7je9h0gv2v55t9mamqrhgja7wypp4yf5d0se3nzyj",
        );
        assert.equal(
            UserSigner.fromWallet(keyFileObjectWithMnemonic, password, 2).getAddress().toBech32(),
            "drt15hpu70r43r3hx9evqmu2z04097ye5t0jgrw3lhxw5rnge0k89nlsh83ydx",
        );

        assert.equal(
            UserSigner.fromWallet(keyFileObjectWithMnemonic, password, 0).getAddress("test").toBech32(),
            "test18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82s4udwpa",
        );
        assert.equal(
            UserSigner.fromWallet(keyFileObjectWithMnemonic, password, 1).getAddress("xdrt").toBech32(),
            "xdrt1tzkwpg0et0s7fp46a7je9h0gv2v55t9mamqrhgja7wypp4yf5d0s7x7qnx",
        );
        assert.equal(
            UserSigner.fromWallet(keyFileObjectWithMnemonic, password, 2).getAddress("ydrt").toBech32(),
            "ydrt15hpu70r43r3hx9evqmu2z04097ye5t0jgrw3lhxw5rnge0k89nls9tf94q",
        );
    });

    it("should throw error when decrypting secret key with keystore-mnemonic file", async function () {
        const userWallet = UserWallet.fromMnemonic({ mnemonic: dummyMnemonic, password: `` });
        const keystoreMnemonic = userWallet.toJSON();

        assert.throws(() => {
            UserWallet.decryptSecretKey(keystoreMnemonic, ``);
        }, `Expected keystore kind to be secretKey, but it was mnemonic.`);
    });
});
