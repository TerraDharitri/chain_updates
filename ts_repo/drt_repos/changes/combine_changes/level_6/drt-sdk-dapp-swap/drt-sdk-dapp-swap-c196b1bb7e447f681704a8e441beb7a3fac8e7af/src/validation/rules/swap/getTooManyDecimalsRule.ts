import { DcdtType } from 'types';
import { canParseAmount } from 'utils';
import { RuleType } from 'validation/types';

export const getTooManyDecimalsRule = (
  token?: DcdtType
): RuleType<string | undefined> => ({
  name: 'decimals',
  message: 'Too many decimals or value too small',
  test: (value) => canParseAmount(String(value), token?.decimals)
});
