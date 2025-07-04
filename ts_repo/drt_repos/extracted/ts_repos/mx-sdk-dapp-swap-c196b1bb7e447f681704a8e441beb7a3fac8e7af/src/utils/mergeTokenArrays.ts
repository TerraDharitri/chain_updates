import { UserEsdtType } from 'types';

export const mergeTokens = (
  tokens: UserEsdtType[] | undefined,
  userTokens: UserEsdtType[] | undefined
): UserEsdtType[] => {
  const tokensMap = new Map<string, UserEsdtType>();

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
      if (userToken.type !== 'FungibleESDT-LP') {
        tokensMap.set(userToken.identifier, userToken);
      }
    }
  });

  return Array.from(tokensMap.values());
};
