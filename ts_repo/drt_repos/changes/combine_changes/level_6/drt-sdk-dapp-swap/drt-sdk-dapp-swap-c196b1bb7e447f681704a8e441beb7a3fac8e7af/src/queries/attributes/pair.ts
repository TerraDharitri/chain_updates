import { dcdtAttributes } from './dcdt';
import { simpleLockAttributes } from './simpleLock';

export const pairAttributes = `
  address
  firstToken {
    ${dcdtAttributes}
  }
  secondToken {
    ${dcdtAttributes}
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
