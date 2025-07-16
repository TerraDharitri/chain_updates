import { assert } from "chai";
import { Address } from "./address";
import { BaseController } from "./baseController";
import { Transaction } from "./transaction";

class TestableBaseController extends BaseController {
    public exposeSetTransactionGasOptions(
        transaction: Transaction,
        options: { gasLimit?: bigint; gasPrice?: bigint },
    ): void {
        this.setTransactionGasOptions(transaction, options);
    }

    public exposeSetVersionAndOptionsForGuardian(transaction: Transaction): void {
        this.setVersionAndOptionsForGuardian(transaction);
    }
}

describe("BaseController Tests", function () {
    it("set correct gasLimit", function () {
        const controller = new TestableBaseController();

        const transaction = new Transaction({
            sender: Address.newFromBech32("drt1yxgvusarlvccy8p30gpl0gywwegd7nqghjqguyn672qwdfdg7y6qqta3dj"),
            receiver: Address.newFromBech32("drt1yxgvusarlvccy8p30gpl0gywwegd7nqghjqguyn672qwdfdg7y6qqta3dj"),
            gasLimit: 0n,
            chainID: "D",
        });

        controller.exposeSetTransactionGasOptions(transaction, { gasLimit: 50000n });
        assert.equal(transaction.gasLimit, 50000n);

        transaction.guardian = Address.newFromBech32("drt1k2s324ww2g0yj38qn2ch2jwctdy8mnfxep94q9arncc6xecg3xaq889n6e");
        transaction.relayer = Address.newFromBech32("drt1kyaqzaprcdnv4luvanah0gfxzzsnpaygsy6pytrexll2urtd05tscswtlq");
        controller.exposeSetTransactionGasOptions(transaction, {});
        assert.equal(transaction.gasLimit, 150000n);
    });

    it("set correct version and options for guarded transactions", function () {
        const controller = new TestableBaseController();

        const transaction = new Transaction({
            sender: Address.newFromBech32("drt1yxgvusarlvccy8p30gpl0gywwegd7nqghjqguyn672qwdfdg7y6qqta3dj"),
            receiver: Address.newFromBech32("drt1yxgvusarlvccy8p30gpl0gywwegd7nqghjqguyn672qwdfdg7y6qqta3dj"),
            gasLimit: 0n,
            chainID: "D",
            version: 0,
            options: 0,
            guardian: Address.newFromBech32("drt1k2s324ww2g0yj38qn2ch2jwctdy8mnfxep94q9arncc6xecg3xaq889n6e"),
        });

        controller.exposeSetVersionAndOptionsForGuardian(transaction);
        assert.equal(transaction.version, 2);
        assert.equal(transaction.options, 2);
    });
});
