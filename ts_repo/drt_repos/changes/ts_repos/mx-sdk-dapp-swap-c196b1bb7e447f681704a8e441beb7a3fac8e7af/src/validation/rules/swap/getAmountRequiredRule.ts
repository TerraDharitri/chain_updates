import { RuleType } from 'validation/types';

export const getAmountRequiredRule = (): RuleType<string | undefined> => ({
  name: 'amount',
  message: 'Amount required',
  test: (value) => Boolean(value)
});
