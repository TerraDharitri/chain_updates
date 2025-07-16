import { assert } from "chai";
import { Address, ErrParseTransactionOutcome, TransactionEvent, TransactionLogs, TransactionOnNetwork } from "../core";
import { b64TopicsToBytes } from "../testutils";
import { SmartContractResult } from "../transactionsOutcomeParsers";
import { TokenManagementTransactionsOutcomeParser } from "./tokenManagementTransactionsOutcomeParser";

describe("test token management transactions outcome parser", () => {
    const parser = new TokenManagementTransactionsOutcomeParser();

    it("should test ensure error", () => {
        const encodedTopics = ["Avk0jZ1kR+l9c76wQQoYcu4hvXPz+jxxTdqQeaCrbX8=", "dGlja2VyIG5hbWUgaXMgbm90IHZhbGlk"];
        const event = new TransactionEvent({
            address: new Address("drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqqqqzlllsd5j0s2"),
            identifier: "signalError",
            topics: b64TopicsToBytes(encodedTopics),
            additionalData: [Buffer.from("QDc1NzM2NTcyMjA2NTcyNzI2Zjcy", "base64")],
        });

        const logs = new TransactionLogs({ events: [event] });
        const transaction = new TransactionOnNetwork({ logs: logs });

        assert.throws(
            () => {
                parser.parseIssueFungible(transaction);
            },
            ErrParseTransactionOutcome,
            "encountered signalError: ticker name is not valid (user error)",
        );
    });

    it("should test parse issue fungible", () => {
        const identifier = "ZZZ-9ee87d";
        const base64Identifier = Buffer.from(identifier).toString("base64");

        const encodedTopics = [base64Identifier, "U0VDT05E", "Wlpa", "RnVuZ2libGVFU0RU", "Ag=="];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "issue",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const logs = new TransactionLogs({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({ logs: logs });

        const outcome = parser.parseIssueFungible(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].tokenIdentifier, identifier);
    });

    it("should test parse issue non fungible", () => {
        const identifier = "NFT-f01d1e";
        const base64Identifier = Buffer.from(identifier).toString("base64");

        let encodedTopics = [
            "TkZULWYwMWQxZQ==",
            "",
            "Y2FuVXBncmFkZQ==",
            "dHJ1ZQ==",
            "Y2FuQWRkU3BlY2lhbFJvbGVz",
            "dHJ1ZQ==",
        ];
        const firstEvent = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "upgradeProperties",
            topics: b64TopicsToBytes(encodedTopics),
        });

        encodedTopics = ["TkZULWYwMWQxZQ==", "", "", "RVNEVFJvbGVCdXJuRm9yQWxs"];
        const secondEvent = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "DCDTSetBurnRoleForAll",
            topics: b64TopicsToBytes(encodedTopics),
        });

        encodedTopics = [base64Identifier, "TkZURVNU", "TkZU", "Tm9uRnVuZ2libGVFU0RU"];
        const thirdEvent = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "issueNonFungible",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const logs = new TransactionLogs({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            events: [firstEvent, secondEvent, thirdEvent],
        });

        const transaction = new TransactionOnNetwork({ logs: logs });

        const outcome = parser.parseIssueNonFungible(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].tokenIdentifier, identifier);
    });

    it("should test parse issue semi fungible", () => {
        const identifier = "SEMIFNG-2c6d9f";
        const base64Identifier = Buffer.from(identifier).toString("base64");

        const encodedTopics = [base64Identifier, "U0VNSQ==", "U0VNSUZORw==", "U2VtaUZ1bmdpYmxlRVNEVA=="];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "issueSemiFungible",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const logs = new TransactionLogs({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({ logs: logs });

        const outcome = parser.parseIssueSemiFungible(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].tokenIdentifier, identifier);
    });

    it("should test parse register meta dcdt", () => {
        const identifier = "METATEST-e05d11";
        const base64Identifier = Buffer.from(identifier).toString("base64");

        const encodedTopics = [base64Identifier, "TUVURVNU", "TUVUQVRFU1Q=", "TWV0YUVTRFQ="];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "registerMetaDCDT",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const logs = new TransactionLogs({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({ logs: logs });

        const outcome = parser.parseRegisterMetaDcdt(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].tokenIdentifier, identifier);
    });

    it("should test parse register and set all roles", () => {
        const firstIdentifier = "LMAO-d9f892";
        const firstBase64Identifier = Buffer.from(firstIdentifier).toString("base64");

        const secondIdentifier = "TST-123456";
        const secondBase64Identifier = Buffer.from(secondIdentifier).toString("base64");

        const roles = ["DCDTRoleLocalMint", "DCDTRoleLocalBurn"];

        let encodedTopics = [firstBase64Identifier, "TE1BTw==", "TE1BTw==", "RnVuZ2libGVFU0RU", "Ag=="];
        const firstEvent = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "registerAndSetAllRoles",
            topics: b64TopicsToBytes(encodedTopics),
        });

        encodedTopics = [secondBase64Identifier, "TE1BTw==", "TE1BTw==", "RnVuZ2libGVFU0RU", "Ag=="];
        const secondEvent = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "registerAndSetAllRoles",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            events: [firstEvent, secondEvent],
        });

        encodedTopics = ["TE1BTy1kOWY4OTI=", "", "", "RENEVFJvbGVMb2NhbE1pbnQ=", "RENEVFJvbGVMb2NhbEJ1cm4="];
        const firstResultEvent = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "DCDTSetRole",
            topics: b64TopicsToBytes(encodedTopics),
        });

        encodedTopics = ["VFNULTEyMzQ1Ng==", "", "", "RENEVFJvbGVMb2NhbE1pbnQ=", "RENEVFJvbGVMb2NhbEJ1cm4="];
        const secondResultEvent = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "DCDTSetRole",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const resultLogs = new TransactionLogs({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            events: [firstResultEvent, secondResultEvent],
        });

        const scResult = new SmartContractResult({
            sender: new Address("drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqqqqzlllsd5j0s2"),
            receiver: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            data: Buffer.from(
                "RENEVFNldFJvbGVANGM0ZDQxNGYyZDY0Mzk2NjM4MzkzMkA0NDQzNDQ1NDUyNmY2YzY1NGM2ZjYzNjE2YzRkNjk2ZTc0QDQ0NDM0NDU0NTI2ZjZjNjU0YzZmNjM2MTZjNDI3NTcyNmU=",
                "base64",
            ),
            logs: resultLogs,
        });

        const transaction = new TransactionOnNetwork({
            smartContractResults: [scResult],
            logs: transactionLogs,
        });

        const outcome = parser.parseRegisterAndSetAllRoles(transaction);
        assert.lengthOf(outcome, 2);

        assert.equal(outcome[0].tokenIdentifier, firstIdentifier);
        assert.deepEqual(outcome[0].roles, roles);

        assert.equal(outcome[1].tokenIdentifier, secondIdentifier);
        assert.deepEqual(outcome[1].roles, roles);
    });

    it("should test parse register set special role", () => {
        const identifier = "METATEST-e05d11";
        const base64Identifier = Buffer.from(identifier).toString("base64");
        const roles = ["DCDTRoleNFTCreate", "DCDTRoleNFTAddQuantity", "DCDTRoleNFTBurn"];

        const encodedTopics = [
            base64Identifier,
            "",
            "",
            "RENEVFJvbGVORlRDcmVhdGU=",
            "RENEVFJvbGVORlRBZGRRdWFudGl0eQ==",
            "RENEVFJvbGVORlRCdXJu",
        ];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "DCDTSetRole",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({
            logs: transactionLogs,
        });

        const outcome = parser.parseSetSpecialRole(transaction);
        assert.lengthOf(outcome, 1);
        assert.deepEqual(
            outcome[0].userAddress,
            new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
        );
        assert.equal(outcome[0].tokenIdentifier, identifier);
        assert.deepEqual(outcome[0].roles, roles);
    });

    it("should test parse nft create", () => {
        const identifier = "NFT-f01d1e";
        const base64Identifier = Buffer.from(identifier).toString("base64");
        const nonce = BigInt(1);
        const initialQuantity = BigInt(1);

        const encodedTopics = [
            base64Identifier,
            "AQ==",
            "AQ==",
            "CAESAgABIuUBCAESCE5GVEZJUlNUGiA8NdfqyxqZpKDMqlN+8MwK4Qn0H2wrQCID5jO/uwcfXCDEEyouUW1ZM3ZKQ3NVcWpNM3hxeGR3VWczemJoVFNMUWZoN0szbW5aWXhyaGNRRFl4RzJDaHR0cHM6Ly9pcGZzLmlvL2lwZnMvUW1ZM3ZKQ3NVcWpNM3hxeGR3VWczemJoVFNMUWZoN0szbW5aWXhyaGNRRFl4Rzo9dGFnczo7bWV0YWRhdGE6UW1SY1A5NGtYcjV6WmpSR3ZpN21KNnVuN0xweFVoWVZSNFI0UnBpY3h6Z1lrdA==",
        ];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "DCDTNFTCreate",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({
            logs: transactionLogs,
        });

        const outcome = parser.parseNftCreate(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].tokenIdentifier, identifier);
        assert.equal(outcome[0].nonce, nonce);
        assert.equal(outcome[0].initialQuantity, initialQuantity);
    });

    it("should test parse local mint", () => {
        const identifier = "AAA-29c4c9";
        const base64Identifier = Buffer.from(identifier).toString("base64");
        const nonce = BigInt(0);
        const mintedSupply = BigInt(100000);

        const encodedTopics = [base64Identifier, "", "AYag"];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "DCDTLocalMint",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({
            logs: transactionLogs,
        });

        const outcome = parser.parseLocalMint(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].userAddress, event.address);
        assert.equal(outcome[0].tokenIdentifier, identifier);
        assert.equal(outcome[0].nonce, nonce);
        assert.equal(outcome[0].mintedSupply, mintedSupply);
    });

    it("should test parse local burn", () => {
        const identifier = "AAA-29c4c9";
        const base64Identifier = Buffer.from(identifier).toString("base64");
        const nonce = BigInt(0);
        const burntSupply = BigInt(100000);

        const encodedTopics = [base64Identifier, "", "AYag"];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "DCDTLocalBurn",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({
            logs: transactionLogs,
        });

        const outcome = parser.parseLocalBurn(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].userAddress, event.address);
        assert.equal(outcome[0].tokenIdentifier, identifier);
        assert.equal(outcome[0].nonce, nonce);
        assert.equal(outcome[0].burntSupply, burntSupply);
    });

    it("should test parse pause", () => {
        const identifier = "AAA-29c4c9";
        const base64Identifier = Buffer.from(identifier).toString("base64");

        const encodedTopics = [base64Identifier];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "DCDTPause",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({
            logs: transactionLogs,
        });

        const outcome = parser.parsePause(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].tokenIdentifier, identifier);
    });

    it("should test parse unpause", () => {
        const identifier = "AAA-29c4c9";
        const base64Identifier = Buffer.from(identifier).toString("base64");

        const encodedTopics = [base64Identifier];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "DCDTUnPause",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({
            logs: transactionLogs,
        });

        const outcome = parser.parseUnpause(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].tokenIdentifier, identifier);
    });

    it("should test parse freeze", () => {
        const identifier = "AAA-29c4c9";
        const base64Identifier = Buffer.from(identifier).toString("base64");
        const nonce = BigInt(0);
        const balance = BigInt(10000000);
        const address = "drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh";

        const encodedTopics = [base64Identifier, "", "mJaA", "OR+TJwep36htO8uz1dDMnyWtDmgP5JnxB9hEt+bqcdU="];
        const event = new TransactionEvent({
            address: new Address("drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqqqqzlllsd5j0s2"),
            identifier: "DCDTFreeze",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            events: [event],
        });

        const scResult = [
            new SmartContractResult({
                sender: new Address("drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqqqqzlllsd5j0s2"),
                receiver: new Address("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
                data: Buffer.from("RENEVEZyZWV6ZUA0MTQxNDEyZDMyMzk2MzM0NjMzOQ=="),
                logs: transactionLogs,
            }),
        ];

        const transaction = new TransactionOnNetwork({
            smartContractResults: scResult,
        });

        const outcome = parser.parseFreeze(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].userAddress, address);
        assert.equal(outcome[0].tokenIdentifier, identifier);
        assert.equal(outcome[0].nonce, nonce);
        assert.equal(outcome[0].balance, balance);
    });

    it("should test parse unfreeze", () => {
        const identifier = "AAA-29c4c9";
        const base64Identifier = Buffer.from(identifier).toString("base64");
        const nonce = BigInt(0);
        const balance = BigInt(10000000);
        const address = "drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh";

        const encodedTopics = [base64Identifier, "", "mJaA", "OR+TJwep36htO8uz1dDMnyWtDmgP5JnxB9hEt+bqcdU="];
        const event = new TransactionEvent({
            address: new Address("drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqqqqzlllsd5j0s2"),
            identifier: "DCDTUnFreeze",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            events: [event],
        });

        const scResult = [
            new SmartContractResult({
                sender: new Address("drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqqqqzlllsd5j0s2"),
                receiver: new Address("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
                data: Buffer.from("RENEVEZyZWV6ZUA0MTQxNDEyZDMyMzk2MzM0NjMzOQ=="),
                logs: transactionLogs,
            }),
        ];

        const transaction = new TransactionOnNetwork({
            smartContractResults: scResult,
        });

        const outcome = parser.parseUnfreeze(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].userAddress, address);
        assert.equal(outcome[0].tokenIdentifier, identifier);
        assert.equal(outcome[0].nonce, nonce);
        assert.equal(outcome[0].balance, balance);
    });

    it("should test parse wipe", () => {
        const identifier = "AAA-29c4c9";
        const base64Identifier = Buffer.from(identifier).toString("base64");
        const nonce = BigInt(0);
        const balance = BigInt(10000000);
        const address = "drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh";

        const encodedTopics = [base64Identifier, "", "mJaA", "OR+TJwep36htO8uz1dDMnyWtDmgP5JnxB9hEt+bqcdU="];
        const event = new TransactionEvent({
            address: new Address("drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqqqqzlllsd5j0s2"),
            identifier: "DCDTWipe",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            events: [event],
        });

        const scResult = [
            new SmartContractResult({
                sender: new Address("drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqqqqzlllsd5j0s2"),
                receiver: new Address("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
                data: Buffer.from("RENEVEZyZWV6ZUA0MTQxNDEyZDMyMzk2MzM0NjMzOQ=="),
                logs: transactionLogs,
            }),
        ];

        const transaction = new TransactionOnNetwork({ smartContractResults: scResult });

        const outcome = parser.parseWipe(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].userAddress, address);
        assert.equal(outcome[0].tokenIdentifier, identifier);
        assert.equal(outcome[0].nonce, nonce);
        assert.equal(outcome[0].balance, balance);
    });

    it("should test parse update attributes", () => {
        const identifier = "NFT-f01d1e";
        const base64Identifier = Buffer.from(identifier).toString("base64");
        const nonce = BigInt(1);
        const attributes = "metadata:ipfsCID/test.json;tags:tag1,tag2";
        const base64Attributes = Buffer.from(attributes).toString("base64");

        const encodedTopics = [base64Identifier, "AQ==", "", base64Attributes];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "DCDTNFTUpdateAttributes",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({
            logs: transactionLogs,
        });

        const outcome = parser.parseUpdateAttributes(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].tokenIdentifier, identifier);
        assert.equal(outcome[0].nonce, nonce);
        assert.equal(Buffer.from(outcome[0].attributes).toString(), attributes);
    });

    it("should test parse add quantity", () => {
        const identifier = "NFT-f01d1e";
        const base64Identifier = Buffer.from(identifier).toString("base64");
        const nonce = BigInt(1);
        const addedQuantity = BigInt(10);

        const encodedTopics = [base64Identifier, "AQ==", "Cg=="];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "DCDTNFTAddQuantity",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({
            logs: transactionLogs,
        });

        const outcome = parser.parseAddQuantity(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].tokenIdentifier, identifier);
        assert.equal(outcome[0].nonce, nonce);
        assert.equal(outcome[0].addedQuantity, addedQuantity);
    });

    it("should test parse burn quantity", () => {
        const identifier = "NFT-f01d1e";
        const base64Identifier = Buffer.from(identifier).toString("base64");
        const nonce = BigInt(1);
        const burntQuantity = BigInt(16);

        const encodedTopics = [base64Identifier, "AQ==", "EA=="];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "DCDTNFTBurn",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({
            logs: transactionLogs,
        });

        const outcome = parser.parseBurnQuantity(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].tokenIdentifier, identifier);
        assert.equal(outcome[0].nonce, nonce);
        assert.equal(outcome[0].burntQuantity, burntQuantity);
    });

    it("should test parse modify royalties", () => {
        const identifier = "TEST-e2b0f9";
        const base64Identifier = Buffer.from(identifier).toString("base64");
        const nonce = BigInt(1);
        const royalties = 20n;

        const encodedTopics = [base64Identifier, "AQ==", "", "FA=="];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "DCDTModifyRoyalties",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({
            logs: transactionLogs,
        });

        const outcome = parser.parseModifyRoyalties(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].tokenIdentifier, identifier);
        assert.equal(outcome[0].nonce, nonce);
        assert.equal(outcome[0].royalties, royalties);
    });

    it("should test parse set new uris", () => {
        const identifier = "TEST-e2b0f9";
        const base64Identifier = Buffer.from(identifier).toString("base64");
        const nonce = BigInt(1);
        const uri = "thisianuri.com";

        const encodedTopics = [base64Identifier, "AQ==", "", "dGhpc2lhbnVyaS5jb20="];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "DCDTSetNewURIs",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({
            logs: transactionLogs,
        });

        const outcome = parser.parseSetNewUris(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].tokenIdentifier, identifier);
        assert.equal(outcome[0].nonce, nonce);
        assert.equal(outcome[0].uri, uri);
    });

    it("should test parse modify creator", () => {
        const identifier = "TEST-e2b0f9";
        const base64Identifier = Buffer.from(identifier).toString("base64");
        const nonce = BigInt(1);

        const encodedTopics = [base64Identifier, "AQ=="];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "DCDTModifyCreator",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({
            logs: transactionLogs,
        });

        const outcome = parser.parseModifyCreator(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].tokenIdentifier, identifier);
        assert.equal(outcome[0].nonce, nonce);
    });

    it("should test parse update metadata", () => {
        const identifier = "TEST-e2b0f9";
        const base64Identifier = Buffer.from(identifier).toString("base64");
        const nonce = BigInt(1);
        const metadata = new Uint8Array(
            Buffer.from(
                "CAUSAgABIlQIARIHVEVTVE5GVBogOR+TJwep36htO8uz1dDMnyWtDmgP5JnxB9hEt+bqcdUgHioIbmV3X2hhc2gyDnRoaXNpYW51cmkuY29tOgkAAAAAAAAAAwUqHgjH9a4DEMf1rgMYx/WuAyDH9a4DKMb1rgMwx/WuAw==",
                "base64",
            ),
        );

        const encodedTopics = [
            base64Identifier,
            "AQ==",
            "",
            "CAUSAgABIlQIARIHVEVTVE5GVBogOR+TJwep36htO8uz1dDMnyWtDmgP5JnxB9hEt+bqcdUgHioIbmV3X2hhc2gyDnRoaXNpYW51cmkuY29tOgkAAAAAAAAAAwUqHgjH9a4DEMf1rgMYx/WuAyDH9a4DKMb1rgMwx/WuAw==",
        ];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "DCDTMetaDataUpdate",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({
            logs: transactionLogs,
        });

        const outcome = parser.parseUpdateMetadata(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].tokenIdentifier, identifier);
        assert.equal(outcome[0].nonce, nonce);
        assert.deepEqual(outcome[0].metadata, metadata);
    });

    it("should test parse recreate metadata", () => {
        const identifier = "TEST-e2b0f9";
        const base64Identifier = Buffer.from(identifier).toString("base64");
        const nonce = BigInt(1);
        const metadata = new Uint8Array(
            Buffer.from(
                "CAUSAgABIlAIARIHVEVTVE5GVBogOR+TJwep36htO8uz1dDMnyWtDmgP5JnxB9hEt+bqcdUgHioSbmV3X2hhc2hfcmVjcmVhdGVkMgA6CQAAAAAAAABkASoeCMj1rgMQyPWuAxjI9a4DIMj1rgMoyPWuAzDI9a4D",
                "base64",
            ),
        );

        const encodedTopics = [
            base64Identifier,
            "AQ==",
            "",
            "CAUSAgABIlAIARIHVEVTVE5GVBogOR+TJwep36htO8uz1dDMnyWtDmgP5JnxB9hEt+bqcdUgHioSbmV3X2hhc2hfcmVjcmVhdGVkMgA6CQAAAAAAAABkASoeCMj1rgMQyPWuAxjI9a4DIMj1rgMoyPWuAzDI9a4D",
        ];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "DCDTMetaDataRecreate",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({
            logs: transactionLogs,
        });

        const outcome = parser.parseMetadataRecreate(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].tokenIdentifier, identifier);
        assert.equal(outcome[0].nonce, nonce);
        assert.deepEqual(outcome[0].metadata, metadata);
    });

    it("should test parse change to dynamic", () => {
        const identifier = "LKXOXNO-503365";
        const base64Identifier = Buffer.from(identifier).toString("base64");
        const tokenName = "LKXOXNO";
        const tokenTicker = "LKXOXNO";
        const tokenType = "DynamicMetaDCDT";

        const encodedTopics = [base64Identifier, "TEtYT1hOTw==", "TEtYT1hOTw==", "RHluYW1pY01ldGFEQ0RU"];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "changeToDynamic",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({
            logs: transactionLogs,
        });

        const outcome = parser.parseChangeTokenToDynamic(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].tokenIdentifier, identifier);
        assert.equal(outcome[0].tokenName, tokenName);
        assert.equal(outcome[0].tickerName, tokenTicker);
        assert.deepEqual(outcome[0].tokenType, tokenType);
    });

    it("should test parse register dynamic", () => {
        const identifier = "TEST-9bbb21";
        const base64Identifier = Buffer.from(identifier).toString("base64");
        const tokenName = "TESTNFT";
        const tokenTicker = "TEST";
        const tokenType = "DynamicNonFungibleDCDT";

        const encodedTopics = [base64Identifier, "VEVTVE5GVA==", "VEVTVA==", "RHluYW1pY05vbkZ1bmdpYmxlRENEVA=="];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "registerDynamic",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({
            logs: transactionLogs,
        });

        const outcome = parser.parseRegisterDynamicToken(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].tokenIdentifier, identifier);
        assert.equal(outcome[0].tokenName, tokenName);
        assert.equal(outcome[0].tokenTicker, tokenTicker);
        assert.deepEqual(outcome[0].tokenType, tokenType);
    });

    it("should test parse register dynamic", () => {
        const identifier = "TEST-9bbb21";
        const base64Identifier = Buffer.from(identifier).toString("base64");
        const tokenName = "TESTNFT";
        const tokenTicker = "TEST";
        const tokenType = "DynamicNonFungibleDCDT";

        const encodedTopics = [base64Identifier, "VEVTVE5GVA==", "VEVTVA==", "RHluYW1pY05vbkZ1bmdpYmxlRENEVA=="];
        const event = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "registerAndSetAllRolesDynamic",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const transactionLogs = new TransactionLogs({
            address: new Address("drt18y0exfc84806smfmeweat5xvnuj66rngpljfnug8mpzt0eh2w82sc0eqzh"),
            events: [event],
        });

        const transaction = new TransactionOnNetwork({
            logs: transactionLogs,
        });

        const outcome = parser.parseRegisterDynamicTokenAndSettingRoles(transaction);
        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].tokenIdentifier, identifier);
        assert.equal(outcome[0].tokenName, tokenName);
        assert.equal(outcome[0].tokenTicker, tokenTicker);
        assert.deepEqual(outcome[0].tokenType, tokenType);
    });
});
