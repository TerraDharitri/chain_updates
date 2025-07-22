import { assert } from "chai";
import { Address, TransactionEvent, TransactionLogs, TransactionOnNetwork } from "../core";
import { b64TopicsToBytes } from "../testutils";
import { SmartContractResult } from "../transactionsOutcomeParsers";
import { DelegationTransactionsOutcomeParser } from "./delegationTransactionsOutcomeParser";

describe("test delegation transactions outcome parser", () => {
    const parser = new DelegationTransactionsOutcomeParser();

    it("should test parseCreateNewDelegationContract ", () => {
        const contractAddress = Address.newFromBech32("drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqqy8llllssrzx6z");
        let encodedTopics = [
            "Q8M8GTdWSAAA",
            "Q8M8GTdWSAAA",
            "AQ==",
            "Q8M8GTdWSAAA",
            "IzMAAAAAAAAAAAAAAAAAAAACMzAAAAAAAAAAABD///8=",
        ];

        const delegateEvent = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "delegate",
            topics: b64TopicsToBytes(encodedTopics),
        });

        encodedTopics = [
            "IzMAAAAAAAAAAAAAAAAAAAACMzAAAAAAAAAAABD///8=",
            "PDXX6ssamaSgzKpTfvDMCuEJ9B9sK0AiA+Yzv7sHH1w=",
        ];
        const scDeployEvent = new TransactionEvent({
            address: new Address("drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqqy8llllssrzx6z"),
            identifier: "SCDeploy",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const logs = new TransactionLogs({ events: [delegateEvent, scDeployEvent] });

        encodedTopics = ["b2g6sUl6beG17FCUIkFwCOTGJjoJJi5SjkP2077e6xA="];
        const scResultEvent = new TransactionEvent({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            identifier: "completedTxEvent",
            topics: b64TopicsToBytes(encodedTopics),
        });

        const scResultLog = new TransactionLogs({
            address: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            events: [scResultEvent],
        });

        const scResult = new SmartContractResult({
            sender: new Address("drt1yvesqqqqqqqqqqqqqqqqqqqqqqqqyvesqqqqqqqqqqqqqqqyllls4jxmwv"),
            receiver: new Address("drt18s6a06ktr2v6fgxv4ffhauxvptssnaqlds45qgsrucemlwc8rawqfgxqg5"),
            data: Buffer.from(
                "QDZmNmJAMjMzMzAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAyMzMzMDAwMDAwMDAwMDAwMDAwMDAxMGZmZmZmZg==",
                "base64",
            ),
            logs: scResultLog,
        });
        const transaction = new TransactionOnNetwork({ smartContractResults: [scResult], logs: logs });

        const outcome = parser.parseCreateNewDelegationContract(transaction);

        assert.lengthOf(outcome, 1);
        assert.equal(outcome[0].contractAddress, contractAddress.toBech32());
    });
});
