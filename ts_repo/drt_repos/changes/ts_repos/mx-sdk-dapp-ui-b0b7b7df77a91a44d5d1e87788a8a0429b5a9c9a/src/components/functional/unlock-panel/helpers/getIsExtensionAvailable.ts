import type { ExtendedWindow, SafeWindowType } from '../unlock-panel.types';

export const getIsExtensionAvailable = () => {
  const safeWindow = typeof window !== 'undefined' ? (window as ExtendedWindow) : ({} as SafeWindowType);

  // Check if either elrondWallet or dharitriWallet exists and has an extensionId
  return Boolean(safeWindow?.elrondWallet) || Boolean(safeWindow?.dharitriWallet);
};
