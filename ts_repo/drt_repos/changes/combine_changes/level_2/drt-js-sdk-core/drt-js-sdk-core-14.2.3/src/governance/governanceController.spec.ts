import { assert } from "chai";
import { Account } from "../accounts";
import { SmartContractQueryResponse } from "../core";
import { Address } from "../core/address";
import { ProxyNetworkProvider } from "../networkProviders";
import { b64TopicsToBytes, MockNetworkProvider } from "../testutils";
import { KeyPair, UserSecretKey } from "../wallet";
import { GovernanceController } from "./governanceController";
import { Vote } from "./resources";

describe("test governance controller", function () {
    const chainID = "D";
    const controller = new GovernanceController({
        chainID: chainID,
        networkProvider: new ProxyNetworkProvider("https://devnet-gateway.dharitri.org"),
    });

    const commitHash = "1db734c0315f9ec422b88f679ccfe3e0197b9d67";
    const governanceAddress = "drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqqqqrlllsf45f4t";

    const aliceBech32 = "drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh";
    const secretKey = UserSecretKey.fromString("2bbcdae7e193924fa0d301e7a12c7defc92a93bc5e587cc968f04fcb86022e1c");
    const keypair = new KeyPair(secretKey);
    const alice = Account.newFromKeypair(keypair);

    it("should create transaction for creating new proposal", async function () {
        const expectedData = `proposal@${Buffer.from(commitHash).toString("hex")}@0a@0f`;

        const transaction = await controller.createTransactionForNewProposal(alice, alice.getNonceThenIncrement(), {
            commitHash: commitHash,
            startVoteEpoch: 10,
            endVoteEpoch: 15,
            nativeTokenAmount: 1000_000000000000000000n,
        });

        assert.equal(transaction.sender.toBech32(), aliceBech32);
        assert.equal(transaction.receiver.toBech32(), governanceAddress);
        assert.equal(transaction.value, 1000_000000000000000000n);
        assert.equal(transaction.chainID, chainID);
        assert.equal(transaction.gasLimit, 50_192_500n);
        assert.equal(transaction.data.toString(), expectedData);
    });

    it("should create transaction for voting", async function () {
        const transaction = await controller.createTransactionForVoting(alice, alice.getNonceThenIncrement(), {
            proposalNonce: 1,
            vote: Vote.YES,
        });

        assert.equal(transaction.sender.toBech32(), aliceBech32);
        assert.equal(transaction.receiver.toBech32(), governanceAddress);
        assert.equal(transaction.value, 0n);
        assert.equal(transaction.chainID, chainID);
        assert.equal(transaction.gasLimit, 5_171_000n);
        assert.equal(transaction.data.toString(), "vote@01@796573");
    });

    it("should create transaction for closing proposal", async function () {
        const transaction = await controller.createTransactionForClosingProposal(alice, alice.getNonceThenIncrement(), {
            proposalNonce: 1,
        });

        assert.equal(transaction.sender.toBech32(), aliceBech32);
        assert.equal(transaction.receiver.toBech32(), governanceAddress);
        assert.equal(transaction.value, 0n);
        assert.equal(transaction.chainID, chainID);
        assert.equal(transaction.gasLimit, 50_074_000n);
        assert.equal(transaction.data.toString(), "closeProposal@01");
    });

    it("should create transaction for clearing ended proposals", async function () {
        const expectedData =
            "clearEndedProposals@391f932707a9dfa86d3bcbb3d5d0cc9f25ad0e680fe499f107d844b7e6ea71d5@2190ce43a3fb31821c317a03f7a08e7650df4c08bc808e127af280e6a5a8f134";

        const transaction = await controller.createTransactionForClearingEndedProposals(
            alice,
            alice.getNonceThenIncrement(),
            {
                proposers: [
                    alice.address,
                    Address.newFromBech32("drt1yxgvusarlvccy8p30gpl0gywwegd7nqghjqguyn672qwdfdg7y6qqta3dj"),
                ],
            },
        );

        assert.equal(transaction.sender.toBech32(), aliceBech32);
        assert.equal(transaction.receiver.toBech32(), governanceAddress);
        assert.equal(transaction.value, 0n);
        assert.equal(transaction.chainID, chainID);
        assert.equal(transaction.gasLimit, 150_273_500n);
        assert.equal(transaction.data.toString(), expectedData);
    });

    it("should create transaction for claiming accumulated fees", async function () {
        const transaction = await controller.createTransactionForClaimingAccumulatedFees(
            alice,
            alice.getNonceThenIncrement(),
            {},
        );

        assert.equal(transaction.sender.toBech32(), aliceBech32);
        assert.equal(transaction.receiver.toBech32(), governanceAddress);
        assert.equal(transaction.value, 0n);
        assert.equal(transaction.chainID, chainID);
        assert.equal(transaction.gasLimit, 1_080_000n);
        assert.equal(transaction.data.toString(), "claimAccumulatedFees");
    });

    it("should create transaction for changing config", async function () {
        const expectedData =
            "changeConfig@31303030303030303030303030303030303030303030@3130303030303030303030303030303030303030@35303030@33303030@36303030";

        const transaction = await controller.createTransactionForChangingConfig(alice, alice.getNonceThenIncrement(), {
            proposalFee: 1000000000000000000000n,
            lastProposalFee: 10000000000000000000n,
            minQuorum: 5000,
            minVetoThreshold: 3000,
            minPassThreshold: 6000,
        });

        assert.equal(transaction.sender.toBech32(), aliceBech32);
        assert.equal(transaction.receiver.toBech32(), governanceAddress);
        assert.equal(transaction.value, 0n);
        assert.equal(transaction.chainID, chainID);
        assert.equal(transaction.gasLimit, 50_237_500n);
        assert.equal(transaction.data.toString(), expectedData);
    });

    it("should get voting power", async function () {
        const provider = new MockNetworkProvider();
        const controller = new GovernanceController({
            chainID: chainID,
            networkProvider: provider,
        });

        provider.mockQueryContractOnFunction(
            "viewVotingPower",
            new SmartContractQueryResponse({
                returnDataParts: [Buffer.from("878678326eac900000", "hex")],
                returnCode: "ok",
                returnMessage: "",
                function: "viewVotingPower",
            }),
        );

        const votingPower = await controller.getVotingPower(alice.address);
        assert.equal(votingPower, 2500_000000000000000000n);
    });

    it("should get config", async function () {
        const provider = new MockNetworkProvider();
        const controller = new GovernanceController({
            chainID: chainID,
            networkProvider: provider,
        });

        provider.mockQueryContractOnFunction(
            "viewConfig",
            new SmartContractQueryResponse({
                returnDataParts: [
                    Buffer.from("1000000000000000000000"),
                    Buffer.from("0.2000"),
                    Buffer.from("0.5000"),
                    Buffer.from("0.3300"),
                    Buffer.from("1"),
                ],
                returnCode: "ok",
                returnMessage: "",
                function: "viewConfig",
            }),
        );

        const config = await controller.getConfig();
        assert.equal(config.proposalFee, 1000_000000000000000000n);
        assert.equal(config.minQuorum, 0.2);
        assert.equal(config.minPassThreshold, 0.5);
        assert.equal(config.minVetoThreshold, 0.33);
        assert.equal(config.lastProposalNonce, 1);
    });

    it("should get proposal", async function () {
        const provider = new MockNetworkProvider();
        const controller = new GovernanceController({
            chainID: chainID,
            networkProvider: provider,
        });

        provider.mockQueryContractOnFunction(
            "viewProposal",
            new SmartContractQueryResponse({
                returnDataParts: b64TopicsToBytes([
                    "NjXJrcXeoAAA",
                    "MWRiNzM0YzAzMTVmOWVjNDIyYjg4ZjY3OWNjZmUzZTAxOTdiOWQ2Nw==",
                    "AQ==",
                    "OR+TJwep36htO8uz1dDMnyWtDmgP5JnxB9hEt+bqcdU=",
                    "NQ==",
                    "Nw==",
                    "",
                    "",
                    "",
                    "",
                    "",
                    "ZmFsc2U=",
                    "ZmFsc2U=",
                ]),
                returnCode: "ok",
                returnMessage: "",
                function: "viewProposal",
            }),
        );

        const proposal = await controller.getProposal(1);
        assert.equal(proposal.cost, 1000_000000000000000000n);
        assert.equal(proposal.commitHash, "1db734c0315f9ec422b88f679ccfe3e0197b9d67");
        assert.equal(proposal.nonce, 1);
        assert.equal(proposal.issuer.toBech32(), aliceBech32);
        assert.equal(proposal.startVoteEpoch, 53);
        assert.equal(proposal.endVoteEpoch, 55);
        assert.equal(proposal.quorumStake, 0n);
        assert.equal(proposal.numYesVotes, 0n);
        assert.equal(proposal.numNoVotes, 0n);
        assert.equal(proposal.numVetoVotes, 0n);
        assert.equal(proposal.numAbstainVotes, 0n);
        assert.equal(proposal.isClosed, false);
        assert.equal(proposal.isPassed, false);
    });
});
