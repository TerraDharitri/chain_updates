import { stringIsFloat } from '@terradharitri/sdk-dapp/utils/validation/stringIsFloat';
import { RuleType } from 'validation/types';

export const getIsValidNumberRule = (): RuleType<string | undefined> => ({
  name: 'isValidNumber',
  message: 'Only digits and one . allowed',
  test: (value) => {
    // run the test only if the input has value
    if (!value || value === '') return true;
    return Boolean(value && stringIsFloat(value));
  }
});
