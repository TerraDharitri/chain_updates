import { homedir } from "os";
import { resolve } from "path";

export const DHARITRI_DIR = resolve(homedir(), ".dharitri");
export const DHARITRI_WALLET_NAME = "wallet.pem";
export const EGLD_NUM_DECIMALS = 18;
export const USER_AGENT = "sdk-dharitri-mcp";
export const MIN_GAS_LIMIT = 50_000n;
