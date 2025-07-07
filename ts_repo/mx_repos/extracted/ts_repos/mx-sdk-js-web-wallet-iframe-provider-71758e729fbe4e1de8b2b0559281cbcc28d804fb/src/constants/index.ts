export const safeWindow = typeof window !== 'undefined' ? window : ({} as any);
export const safeDocument =
  typeof document !== 'undefined' ? document : ({} as any);

export const iframeWindowReadyEvent = 'iframeWindowReady';

export enum IframeLoginTypes {
  metamask = 'metamask',
  passkey = 'passkey'
}
