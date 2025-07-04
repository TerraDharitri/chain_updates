import { collectionAttributes } from './collection';

export const simpleLockAttributes = `
  address
  lockedToken {
    ${collectionAttributes}
  }
  lpProxyToken {
    ${collectionAttributes}
  }
  farmProxyToken {
    ${collectionAttributes}
  }
  intermediatedPairs
  intermediatedFarms
`;
