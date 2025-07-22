import { assert } from "chai";
import { Abi, Code, U32Value } from "../abi";
import { Address, Err, Token, TokenTransfer, TransactionsFactoryConfig } from "../core";
import { loadAbiRegistry, loadContractCode } from "../testutils/utils";
import { SmartContractTransactionsFactory } from "./smartContractTransactionsFactory";

describe("test smart contract transactions factory", function () {
    const config = new TransactionsFactoryConfig({ chainID: "D" });
    let factory: SmartContractTransactionsFactory;
    let abiAwareFactory: SmartContractTransactionsFactory;
    let bytecode: Code;
    let abi: Abi;

    before(async function () {
        factory = new SmartContractTransactionsFactory({
            config: config,
        });

        bytecode = await loadContractCode("src/testdata/adder.wasm");
        abi = await loadAbiRegistry("src/testdata/adder.abi.json");

        abiAwareFactory = new SmartContractTransactionsFactory({
            config: config,
            abi: abi,
        });
    });

    it("should throw error when args are not of type 'TypedValue'", async function () {
        const sender = Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
        const gasLimit = 6000000n;
        const args = [0];

        assert.throws(
            () =>
                factory.createTransactionForDeploy(sender, {
                    bytecode: bytecode.valueOf(),
                    gasLimit: gasLimit,
                    arguments: args,
                }),
            Err,
            "Can't convert args to TypedValues",
        );
    });

    it("should create 'Transaction' for deploy", async function () {
        const sender = Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
        const gasLimit = 6000000n;
        const args = [new U32Value(1)];

        const transaction = factory.createTransactionForDeploy(sender, {
            bytecode: bytecode.valueOf(),
            gasLimit: gasLimit,
            arguments: args,
        });

        const transactionAbiAware = abiAwareFactory.createTransactionForDeploy(sender, {
            bytecode: bytecode.valueOf(),
            gasLimit: gasLimit,
            arguments: args,
        });

        assert.deepEqual(
            transaction.sender,
            Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
        );
        assert.deepEqual(
            transaction.receiver,
            Address.newFromBech32("drt1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq85hk5z"),
        );
        assert.deepEqual(transaction.data, Buffer.from(`${bytecode}@0500@0504@01`));
        assert.equal(transaction.gasLimit.valueOf(), gasLimit);
        assert.equal(transaction.value, 0n);

        assert.deepEqual(transaction, transactionAbiAware);
    });

    it("should create 'Transaction' for execute without transfer", async function () {
        const sender = Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
        const contract = Address.newFromBech32("drt1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fq4h4xut");
        const func = "add";
        const gasLimit = 6000000n;
        const args = [new U32Value(7)];

        const transaction = factory.createTransactionForExecute(sender, {
            contract: contract,
            function: func,
            gasLimit: gasLimit,
            arguments: args,
        });

        const transactionAbiAware = abiAwareFactory.createTransactionForExecute(sender, {
            contract: contract,
            function: func,
            gasLimit: gasLimit,
            arguments: args,
        });

        assert.deepEqual(
            transaction.sender,
            Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
        );
        assert.deepEqual(
            transaction.receiver,
            Address.newFromBech32("drt1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fq4h4xut"),
        );
        assert.deepEqual(transaction.data, Buffer.from("add@07"));
        assert.equal(transaction.gasLimit, gasLimit);
        assert.equal(transaction.value, 0n);

        assert.deepEqual(transaction, transactionAbiAware);
    });

    it("should create 'Transaction' for execute and transfer native token", async function () {
        const sender = Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
        const contract = Address.newFromBech32("drt1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fq4h4xut");
        const func = "add";
        const gasLimit = 6000000n;
        const rewaAmount = 1000000000000000000n;

        const transaction = factory.createTransactionForExecute(sender, {
            contract: contract,
            function: func,
            gasLimit: gasLimit,
            arguments: [new U32Value(7)],
            nativeTransferAmount: rewaAmount,
        });

        const transactionAbiAware = abiAwareFactory.createTransactionForExecute(sender, {
            contract: contract,
            function: func,
            gasLimit: gasLimit,
            arguments: [7],
            nativeTransferAmount: rewaAmount,
        });

        assert.deepEqual(
            transaction.sender,
            Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
        );
        assert.deepEqual(
            transaction.receiver,
            Address.newFromBech32("drt1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fq4h4xut"),
        );
        assert.deepEqual(transaction.data, Buffer.from("add@07"));
        assert.equal(transaction.gasLimit, gasLimit);
        assert.equal(transaction.value, 1000000000000000000n);

        assert.deepEqual(transaction, transactionAbiAware);
    });

    it("should create 'Transaction' for execute and transfer single dcdt", async function () {
        const sender = Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
        const contract = Address.newFromBech32("drt1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fq4h4xut");
        const func = "add";
        const gasLimit = 6000000n;
        const args = [new U32Value(7)];
        const token = new Token({ identifier: "FOO-6ce17b", nonce: 0n });
        const transfer = new TokenTransfer({ token, amount: 10n });

        const transaction = factory.createTransactionForExecute(sender, {
            contract: contract,
            function: func,
            gasLimit: gasLimit,
            arguments: args,
            tokenTransfers: [transfer],
        });

        const transactionAbiAware = abiAwareFactory.createTransactionForExecute(sender, {
            contract: contract,
            function: func,
            gasLimit: gasLimit,
            arguments: args,
            tokenTransfers: [transfer],
        });

        assert.deepEqual(
            transaction.sender,
            Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
        );
        assert.deepEqual(
            transaction.receiver,
            Address.newFromBech32("drt1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fq4h4xut"),
        );
        assert.deepEqual(transaction.data, Buffer.from("DCDTTransfer@464f4f2d366365313762@0a@616464@07"));
        assert.equal(transaction.gasLimit, gasLimit);
        assert.equal(transaction.value, 0n);

        assert.deepEqual(transaction, transactionAbiAware);
    });

    it("should create 'Transaction' for execute and transfer multiple dcdts", async function () {
        const sender = Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
        const contract = Address.newFromBech32("drt1qqqqqqqqqqqqqpgqak8zt22wl2ph4tswtyc39namqx6ysa2sd8ssg6vu30");
        const func = "add";
        const gasLimit = 6000000n;
        const args = [new U32Value(7)];

        const fooToken = new Token({ identifier: "FOO-6ce17b", nonce: 0n });
        const fooTransfer = new TokenTransfer({ token: fooToken, amount: 10n });
        const barToken = new Token({ identifier: "BAR-5bc08f", nonce: 0n });
        const barTransfer = new TokenTransfer({ token: barToken, amount: 3140n });

        const transaction = factory.createTransactionForExecute(sender, {
            contract: contract,
            function: func,
            gasLimit: gasLimit,
            arguments: args,
            tokenTransfers: [fooTransfer, barTransfer],
        });

        const transactionAbiAware = abiAwareFactory.createTransactionForExecute(sender, {
            contract: contract,
            function: func,
            gasLimit: gasLimit,
            arguments: args,
            tokenTransfers: [fooTransfer, barTransfer],
        });

        assert.deepEqual(
            transaction.sender,
            Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
        );
        assert.deepEqual(
            transaction.receiver,
            Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
        );

        assert.deepEqual(
            transaction.data,
            Buffer.from(
                "MultiDCDTNFTTransfer@00000000000000000500ed8e25a94efa837aae0e593112cfbb01b448755069e1@02@464f4f2d366365313762@@0a@4241522d356263303866@@0c44@616464@07",
            ),
        );

        assert.equal(transaction.gasLimit, gasLimit);
        assert.equal(transaction.value, 0n);

        assert.deepEqual(transaction, transactionAbiAware);
    });

    it("should create 'Transaction' for execute and transfer single nft", async function () {
        const sender = Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
        const contract = Address.newFromBech32("drt1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fq4h4xut");
        const func = "add";
        const gasLimit = 6000000n;
        const args = [new U32Value(7)];

        const token = new Token({ identifier: "NFT-123456", nonce: 1n });
        const transfer = new TokenTransfer({ token, amount: 1n });

        const transaction = factory.createTransactionForExecute(sender, {
            contract: contract,
            function: func,
            gasLimit: gasLimit,
            arguments: args,
            tokenTransfers: [transfer],
        });

        const transactionAbiAware = abiAwareFactory.createTransactionForExecute(sender, {
            contract: contract,
            function: func,
            gasLimit: gasLimit,
            arguments: args,
            tokenTransfers: [transfer],
        });

        assert.deepEqual(
            transaction.sender,
            Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
        );
        assert.deepEqual(
            transaction.receiver,
            Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
        );

        assert.isDefined(transaction.data);
        assert.deepEqual(
            transaction.data,
            Buffer.from(
                "DCDTNFTTransfer@4e46542d313233343536@01@01@00000000000000000500b9353fe8407f87310c87e12fa1ac807f0485da39d152@616464@07",
            ),
        );

        assert.equal(transaction.gasLimit, gasLimit);
        assert.equal(transaction.value, 0n);

        assert.deepEqual(transaction, transactionAbiAware);
    });

    it("should create 'Transaction' for execute and transfer multiple nfts", async function () {
        const sender = Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
        const contract = Address.newFromBech32("drt1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fq4h4xut");
        const func = "add";
        const gasLimit = 6000000n;
        const args = [new U32Value(7)];

        const firstToken = new Token({ identifier: "NFT-123456", nonce: 1n });
        const firstTransfer = new TokenTransfer({ token: firstToken, amount: 1n });
        const secondToken = new Token({ identifier: "NFT-123456", nonce: 42n });
        const secondTransfer = new TokenTransfer({ token: secondToken, amount: 1n });

        const transaction = factory.createTransactionForExecute(sender, {
            contract: contract,
            function: func,
            gasLimit: gasLimit,
            arguments: args,
            tokenTransfers: [firstTransfer, secondTransfer],
        });

        const transactionAbiAware = abiAwareFactory.createTransactionForExecute(sender, {
            contract: contract,
            function: func,
            gasLimit: gasLimit,
            arguments: args,
            tokenTransfers: [firstTransfer, secondTransfer],
        });

        assert.deepEqual(
            transaction.sender,
            Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
        );
        assert.deepEqual(
            transaction.receiver,
            Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
        );

        assert.isDefined(transaction.data);
        assert.deepEqual(
            transaction.data,
            Buffer.from(
                "MultiDCDTNFTTransfer@00000000000000000500b9353fe8407f87310c87e12fa1ac807f0485da39d152@02@4e46542d313233343536@01@01@4e46542d313233343536@2a@01@616464@07",
            ),
        );

        assert.equal(transaction.gasLimit, gasLimit);
        assert.equal(transaction.value, 0n);

        assert.deepEqual(transaction, transactionAbiAware);
    });

    it("should create 'Transaction' for execute and transfer native and nfts", async function () {
        const sender = Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
        const contract = Address.newFromBech32("drt1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fq4h4xut");
        const func = "add";
        const gasLimit = 6000000n;
        const args = [new U32Value(7)];

        const firstToken = new Token({ identifier: "NFT-123456", nonce: 1n });
        const firstTransfer = new TokenTransfer({ token: firstToken, amount: 1n });
        const secondToken = new Token({ identifier: "NFT-123456", nonce: 42n });
        const secondTransfer = new TokenTransfer({ token: secondToken, amount: 1n });

        const transaction = factory.createTransactionForExecute(sender, {
            contract: contract,
            function: func,
            gasLimit: gasLimit,
            arguments: args,
            nativeTransferAmount: 1000000000000000000n,
            tokenTransfers: [firstTransfer, secondTransfer],
        });

        const transactionAbiAware = abiAwareFactory.createTransactionForExecute(sender, {
            contract: contract,
            function: func,
            gasLimit: gasLimit,
            arguments: args,
            nativeTransferAmount: 1000000000000000000n,
            tokenTransfers: [firstTransfer, secondTransfer],
        });

        assert.deepEqual(
            transaction.sender,
            Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
        );
        assert.deepEqual(
            transaction.receiver,
            Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
        );

        assert.isDefined(transaction.data);
        assert.deepEqual(
            transaction.data,
            Buffer.from(
                "MultiDCDTNFTTransfer@00000000000000000500b9353fe8407f87310c87e12fa1ac807f0485da39d152@03@4e46542d313233343536@01@01@4e46542d313233343536@2a@01@524557412d303030303030@@0de0b6b3a7640000@616464@07",
            ),
        );

        assert.equal(transaction.gasLimit, gasLimit);
        assert.equal(transaction.value, 0n);

        assert.deepEqual(transaction, transactionAbiAware);
    });

    it("should create 'Transaction' for upgrade", async function () {
        const sender = Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
        const contract = Address.newFromBech32("drt1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fq4h4xut");
        const gasLimit = 6000000n;
        const args = [new U32Value(7)];

        const transaction = factory.createTransactionForUpgrade(sender, {
            contract: contract,
            bytecode: bytecode.valueOf(),
            gasLimit: gasLimit,
            arguments: args,
        });

        const transactionAbiAware = abiAwareFactory.createTransactionForUpgrade(sender, {
            contract: contract,
            bytecode: bytecode.valueOf(),
            gasLimit: gasLimit,
            arguments: args,
        });

        assert.deepEqual(
            transaction.sender,
            Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
        );
        assert.deepEqual(
            transaction.receiver,
            Address.newFromBech32("drt1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fq4h4xut"),
        );
        assert.deepEqual(transaction.data!, Buffer.from(`upgradeContract@${bytecode}@0504@07`));
        assert.equal(transaction.gasLimit, gasLimit);
        assert.equal(transaction.value, 0n);

        assert.deepEqual(transaction, transactionAbiAware);
    });

    it("should create 'Transaction' for upgrade, when ABI is available (with fallbacks)", async function () {
        const abi = Abi.create({
            upgradeConstructor: {
                inputs: [
                    {
                        type: "u32",
                    },
                    {
                        type: "u32",
                    },
                    {
                        type: "u32",
                    },
                ],
            },
            endpoints: [
                {
                    name: "upgrade",
                    inputs: [
                        {
                            type: "u32",
                        },
                        {
                            type: "u32",
                        },
                    ],
                },
            ],
            constructor: {
                inputs: [
                    {
                        type: "u32",
                    },
                ],
            },
        });

        const factory = new SmartContractTransactionsFactory({
            config: config,
            abi: abi,
        });

        const bytecode = Buffer.from("abba", "hex");
        const sender = Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
        const receiver = Address.newFromBech32("drt1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fq4h4xut");
        const gasLimit = 6000000n;

        // By default, use the upgrade constructor.
        const tx1 = factory.createTransactionForUpgrade(sender, {
            contract: receiver,
            bytecode: bytecode,
            gasLimit: gasLimit,
            arguments: [42, 42, 42],
        });

        assert.equal(Buffer.from(tx1.data!).toString(), `upgradeContract@abba@0504@2a@2a@2a`);

        // Fallback to the "upgrade" endpoint.
        (<any>abi).upgradeConstructorDefinition = undefined;

        const tx2 = factory.createTransactionForUpgrade(sender, {
            contract: receiver,
            bytecode: bytecode,
            gasLimit: gasLimit,
            arguments: [42, 42],
        });

        assert.equal(Buffer.from(tx2.data!).toString(), `upgradeContract@abba@0504@2a@2a`);

        // Fallback to the constructor.
        (<any>abi).endpoints.length = 0;

        const tx3 = factory.createTransactionForUpgrade(sender, {
            contract: receiver,
            bytecode: bytecode,
            gasLimit: gasLimit,
            arguments: [42],
        });

        assert.equal(Buffer.from(tx3.data!).toString(), `upgradeContract@abba@0504@2a`);

        // No fallbacks.
        (<any>abi).constructorDefinition = undefined;

        assert.throws(
            () =>
                factory.createTransactionForUpgrade(sender, {
                    contract: receiver,
                    bytecode: bytecode,
                    gasLimit: gasLimit,
                    arguments: [42],
                }),
            "Can't convert args to TypedValues",
        );
    });

    it("should create 'Transaction' for claiming developer rewards", async function () {
        const sender = Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
        const contract = Address.newFromBech32("drt1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fq4h4xut");

        const transaction = factory.createTransactionForClaimingDeveloperRewards({
            sender: sender,
            contract: contract,
        });

        assert.deepEqual(
            transaction.sender,
            Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
        );
        assert.deepEqual(
            transaction.receiver,
            Address.newFromBech32("drt1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fq4h4xut"),
        );
        assert.equal(Buffer.from(transaction.data).toString(), "ClaimDeveloperRewards");
        assert.equal(transaction.gasLimit, 6000000n);
        assert.equal(transaction.value, 0n);
    });

    it("should create 'Transaction' for changing owner address", async function () {
        const sender = Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh");
        const contract = Address.newFromBech32("drt1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fq4h4xut");
        const newOwner = Address.newFromBech32("drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c");

        const transaction = factory.createTransactionForChangingOwnerAddress({
            sender: sender,
            contract: contract,
            newOwner: newOwner,
        });

        assert.deepEqual(
            transaction.sender,
            Address.newFromBech32("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
        );
        assert.deepEqual(
            transaction.receiver,
            Address.newFromBech32("drt1qqqqqqqqqqqqqpgqhy6nl6zq07rnzry8uyh6rtyq0uzgtk3e69fq4h4xut"),
        );
        assert.equal(
            Buffer.from(transaction.data).toString(),
            "ChangeOwnerAddress@8049d639e5a6980d1cd2392abcce41029cda74a1563523a202f09641cc2618f8",
        );
        assert.equal(transaction.gasLimit, 6000000n);
        assert.equal(transaction.value, 0n);
    });
});
