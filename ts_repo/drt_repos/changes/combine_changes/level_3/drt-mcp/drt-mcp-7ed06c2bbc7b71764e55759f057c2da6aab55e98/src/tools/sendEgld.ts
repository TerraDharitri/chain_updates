import { CallToolResult } from "@modelcontextprotocol/sdk/types.js";
import { Account, Address } from "@terradharitri/sdk-core";
import { z } from "zod";
import {
  denominateRewaValue,
  getEntrypoint,
  getExplorerUrl,
  loadNetworkFromEnv,
  loadPemWalletFromEnv,
} from "./utils.js";

export async function sendRewa(
  amount: string,
  receiver: string
): Promise<CallToolResult> {
  const denominated = denominateRewaValue(amount);
  const pem = loadPemWalletFromEnv();
  const receiverAddress = Address.newFromBech32(receiver);
  const account = new Account(pem.secretKey);

  const network = loadNetworkFromEnv();
  const entrypoint = getEntrypoint(network);

  const accountOnNetwork = await entrypoint
    .createNetworkProvider()
    .getAccount(account.address);

  if (denominated > accountOnNetwork.balance) {
    throw new Error("Not enough REWA balance");
  }

  account.nonce = accountOnNetwork.nonce;

  const controller = entrypoint.createTransfersController();
  const transaction = await controller.createTransactionForTransfer(
    account,
    account.nonce,
    {
      receiver: receiverAddress,
      nativeAmount: denominated,
    }
  );

  const hash = await entrypoint.sendTransaction(transaction);
  const explorer = getExplorerUrl(network);
  return {
    content: [
      {
        type: "text",
        text: `${amount} REWA have been sent to ${receiverAddress.toBech32()}. Check out the transaction here: ${explorer}/transactions/${hash}`,
      },
    ],
  };
}

export const sendRewaToolName = "send-rewa";
export const sendRewaToolDescription =
  "Create a move balance transaction and send it. Will send REWA using the wallet set in the env to the specified receiver.";
export const sendRewaParamScheme = {
  amount: z
    .string()
    .describe(
      "The amount of REWA to send. This amount will then be denominated (1 REWA=1000000000000000000)"
    ),
  receiver: z.string().describe("The bech32 address of the receiver (drt1...)"),
};
