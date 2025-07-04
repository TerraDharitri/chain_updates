import BigNumber from 'bignumber.js';
import { EGLD_IDENTIFIER } from 'constants/general';
import { EsdtType, UserEsdtType } from 'types';

export const getSortedTokensByUsdValue = ({
  tokens,
  wrappedEgld
}: {
  tokens: UserEsdtType[];
  wrappedEgld?: EsdtType;
}) => {
  return tokens.sort((esdt1, esdt2) => {
    if (esdt1.identifier === EGLD_IDENTIFIER) return -1;
    if (esdt2.identifier === EGLD_IDENTIFIER) return 1;
    if (esdt1.identifier === wrappedEgld?.identifier) return -1;
    if (esdt2.identifier === wrappedEgld?.identifier) return 1;
    return new BigNumber(esdt2.valueUSD).minus(esdt1.valueUSD).toNumber();
  });
};
