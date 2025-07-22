import { Operation, OptionalOperation } from "./operation";

// WalletConnect Namespace for DharitrI
export const WALLETCONNECT_DHARITRI_NAMESPACE = "drt";
// WalletConnect default methods for DharitrI
export const WALLETCONNECT_DHARITRI_METHODS = Object.values(Operation);
// WalletConnect optional methods for DharitrI
export const WALLETCONNECT_DHARITRI_OPTIONAL_METHODS =
  Object.values(OptionalOperation);
// Delay the sign login token action for 500ms to allow the UI to update properly
export const WALLETCONNECT_SIGN_LOGIN_DELAY = 500;
