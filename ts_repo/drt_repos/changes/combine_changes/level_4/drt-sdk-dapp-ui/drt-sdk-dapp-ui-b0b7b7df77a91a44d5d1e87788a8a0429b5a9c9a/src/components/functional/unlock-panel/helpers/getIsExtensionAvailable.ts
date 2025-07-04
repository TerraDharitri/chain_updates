import type { ExtendedWindow, SafeWindowType } from '../unlock-panel.types';

export const getIsExtensionAvailable = () => {
  const safeWindow = typeof window !== 'undefined' ? (window as ExtendedWindow) : ({} as SafeWindowType);

  // Check if either numbatWallet or dharitriWallet exists and has an extensionId
  return Boolean(safeWindow?.numbatWallet) || Boolean(safeWindow?.dharitriWallet);
};
