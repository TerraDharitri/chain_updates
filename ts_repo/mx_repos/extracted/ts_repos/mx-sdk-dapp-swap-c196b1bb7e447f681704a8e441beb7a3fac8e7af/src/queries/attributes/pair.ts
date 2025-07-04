import { esdtAttributes } from './esdt';
import { simpleLockAttributes } from './simpleLock';

export const pairAttributes = `
  address
  firstToken {
    ${esdtAttributes}
  }
  secondToken {
    ${esdtAttributes}
  }
  type
  feesAPR
  totalFeePercent
  specialFeePercent

  lockedTokensInfo {
    lockingSC {
      ${simpleLockAttributes}
    }
    unlockEpoch
  }
`;
