import { Address } from '@terradharitri/sdk-core';

export const addressIsErd = (address: string) => {
  try {
    const addressObj = Address.newFromBech32(address);
    return addressObj.toBech32().startsWith('drt');
  } catch (error) {
    return false;
  }
};
