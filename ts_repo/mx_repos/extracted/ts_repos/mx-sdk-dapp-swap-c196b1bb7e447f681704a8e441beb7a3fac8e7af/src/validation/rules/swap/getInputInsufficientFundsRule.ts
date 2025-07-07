import BigNumber from 'bignumber.js';
import { parseAmount } from 'lib';
import { EsdtType } from 'types';
import { RuleType } from 'validation/types';

export const getInputInsufficientFundsRule = (
  token?: EsdtType
): RuleType<string | undefined> => ({
  name: 'funds',
  message: 'Insufficient funds',
  test: (amount) => {
    if (amount && token && token.balance != null) {
      const parsedAmount = parseAmount(amount.toString(), token.decimals);
      const bnAmount = new BigNumber(parsedAmount);
      const bnBalance = new BigNumber(token.balance ?? '0');

      return bnBalance.comparedTo(bnAmount) >= 0;
    }
    return true;
  }
});
