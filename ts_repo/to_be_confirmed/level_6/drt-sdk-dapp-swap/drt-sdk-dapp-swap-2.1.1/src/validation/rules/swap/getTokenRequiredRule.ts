import { RuleType } from 'validation/types';

export const getTokenRequiredRule = (
  tokenId?: string
): RuleType<string | undefined> => ({
  name: 'token',
  message: 'Token required',
  test: () => Boolean(tokenId)
});
