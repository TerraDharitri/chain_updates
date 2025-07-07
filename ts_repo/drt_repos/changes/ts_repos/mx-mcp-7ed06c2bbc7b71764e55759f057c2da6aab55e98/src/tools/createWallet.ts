import { CallToolResult } from "@modelcontextprotocol/sdk/types.js";
import { UserPem, UserSecretKey } from "@terradharitri/sdk-core/out/index.js";
import * as fs from "fs";
import { resolve } from "path";
import { DHARITRI_DIR, DHARITRI_WALLET_NAME } from "./constants.js";

export async function createWallet(): Promise<CallToolResult> {
  if (!fs.existsSync(DHARITRI_DIR)) {
    fs.mkdirSync(DHARITRI_DIR, { recursive: true });
  }

  const walletPath = resolve(DHARITRI_DIR, DHARITRI_WALLET_NAME);

  if (fs.existsSync(walletPath)) {
    return {
      content: [
        {
          type: "text",
          text: `A wallet exists at location ${walletPath}. Will not overwrite it.`,
        },
      ],
    };
  }

  const secretKey = UserSecretKey.generate();
  const address = secretKey.generatePublicKey().toAddress();
  const pem = new UserPem(address.toBech32(), secretKey);

  pem.save(walletPath);
  fs.chmodSync(walletPath, 0o444);
  return {
    content: [
      {
        type: "text",
        text: `A wallet has been created and saved as a PEM file at: ${walletPath}. PEM files ARE NOT SECURE.
If you want to further use the generated wallet, make sure to fund it first and set the absolute path in the config file under the "MVX_WALLET" environment variable.`,
      },
    ],
  };
}

export const createWalletToolName = "create-wallet";
export const createWalletToolDescription =
  "Create a new wallet and save it as a PEM file. PEM file ARE NOT SECURE. If a wallet already exists, will abort operation.";
