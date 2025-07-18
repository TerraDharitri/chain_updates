import BigNumber from 'bignumber.js';
import { REWA_IDENTIFIER } from 'constants/general';
import { DcdtType, UserDcdtType } from 'types';

export const getSortedTokensByUsdValue = ({
  tokens,
  wrappedRewa
}: {
  tokens: UserDcdtType[];
  wrappedRewa?: DcdtType;
}) => {
  return tokens.sort((dcdt1, dcdt2) => {
    if (dcdt1.identifier === REWA_IDENTIFIER) return -1;
    if (dcdt2.identifier === REWA_IDENTIFIER) return 1;
    if (dcdt1.identifier === wrappedRewa?.identifier) return -1;
    if (dcdt2.identifier === wrappedRewa?.identifier) return 1;
    return new BigNumber(dcdt2.valueUSD).minus(dcdt1.valueUSD).toNumber();
  });
};
