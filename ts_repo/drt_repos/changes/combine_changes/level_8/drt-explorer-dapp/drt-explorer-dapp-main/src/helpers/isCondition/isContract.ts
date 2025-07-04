import { Address } from '@terradharitri/sdk-core/out';

export const isContract = (initiator: string | undefined) => {
  if (!initiator) {
    return false;
  }

  try {
    return new Address(initiator).isSmartContract();
  } catch {
    return false;
  }
};
