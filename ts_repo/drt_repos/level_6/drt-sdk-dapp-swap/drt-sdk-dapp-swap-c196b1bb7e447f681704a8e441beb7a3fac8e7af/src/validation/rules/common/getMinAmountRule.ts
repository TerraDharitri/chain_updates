import { RuleType } from 'validation/types';

type MinAmountRuleType = {
  minAcceptedAmount?: number;
};

export const getMinAmountRule = ({
  minAcceptedAmount
}: MinAmountRuleType): RuleType<string | undefined> => ({
  name: 'minAmount',
  message: `Minimum amount: ${minAcceptedAmount}`,
  test: (amount) => {
    if (amount && minAcceptedAmount) {
      return parseFloat(amount) >= minAcceptedAmount;
    }
    return true;
  }
});
