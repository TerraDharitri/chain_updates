export * from './windowProviderConstants';

export const safeWindow = typeof window !== 'undefined' ? window : ({} as any);
export const safeDocument =
  typeof document !== 'undefined' ? document : ({} as any);
