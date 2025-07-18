import { CallToolResult } from "@modelcontextprotocol/sdk/types.js";
import { Account, Address } from "@terradharitri/sdk-core/out/index.js";
import { z } from "zod";
import { MIN_GAS_LIMIT } from "./constants.js";
import {
  denominateRewaValue,
  getEntrypoint,
  loadNetworkFromEnv,
  loadPemWalletFromEnv,
} from "./utils.js";

export async function sendRewaToMultipleReceivers(
  amount: string,
  receivers: string[]
): Promise<CallToolResult> {
  const denominated = denominateRewaValue(amount);
  const pem = loadPemWalletFromEnv();

  const account = new Account(pem.secretKey);

  const network = loadNetworkFromEnv();
  const entrypoint = getEntrypoint(network);

  const accountOnNetwork = await entrypoint
    .createNetworkProvider()
    .getAccount(account.address);

  const requiredBalance =
    denominated * BigInt(receivers.length) +
    MIN_GAS_LIMIT * BigInt(receivers.length);

  if (requiredBalance > accountOnNetwork.balance) {
    throw new Error("Not enough REWA balance");
  }

  let hashes: string[] = [];

  const controller = entrypoint.createTransfersController();
  account.nonce = accountOnNetwork.nonce;

  receivers.forEach(async (receiver) => {
    const receiverAddress = Address.newFromBech32(receiver);

    const transaction = await controller.createTransactionForTransfer(
      account,
      account.getNonceThenIncrement(),
      {
        receiver: receiverAddress,
        nativeAmount: denominated,
      }
    );

    const hash = await entrypoint.sendTransaction(transaction);
    hashes.push(hash);
  });

  return {
    content: [
      {
        type: "text",
        text: `${amount} REWA has been sent to each receiver. Check out the transactions hashes here: ${hashes.join(
          ", "
        )}`,
      },
    ],
  };
}

export const sendRewaToMultipleReceiversToolName =
  "send-rewa-to-multiple-receivers";
export const sendRewaToMultipleReceiversToolDescription =
  "Create move balance transactions and send them. Will send REWA using the wallet set in the env to each specified receiver.";
export const sendRewaToMultipleReceiversParamScheme = {
  amount: z
    .string()
    .describe(
      "The amount of REWA to send. This amount will then be denominated (1 REWA=1000000000000000000)"
    ),
  receivers: z
    .array(z.string())
    .describe("An array of bech32 addresses of the receivers (drt1...)"),
};
