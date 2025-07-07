import { UserDcdtType } from 'types';

export const mergeTokens = (
  tokens: UserDcdtType[] | undefined,
  userTokens: UserDcdtType[] | undefined
): UserDcdtType[] => {
  const tokensMap = new Map<string, UserDcdtType>();

  tokens?.forEach((token) => {
    tokensMap.set(token.identifier, {
      ...token,
      balance: '0', // Default balance
      valueUSD: '0' // Default valueUSD
    });
  });

  userTokens?.forEach((userToken) => {
    const existingToken = tokensMap.get(userToken.identifier);

    if (existingToken) {
      tokensMap.set(userToken.identifier, {
        ...existingToken,
        balance: userToken.balance ?? '0',
        valueUSD: userToken.valueUSD ?? '0'
      });
    } else {
      // Filter out LP tokens
      if (userToken.type !== 'FungibleDCDT-LP') {
        tokensMap.set(userToken.identifier, userToken);
      }
    }
  });

  return Array.from(tokensMap.values());
};
