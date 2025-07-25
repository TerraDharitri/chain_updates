import { DECIMALS } from '@terradharitri/sdk-dapp/constants';
import { REWA_IDENTIFIER } from 'constants/general';
import { DcdtType, PairType } from 'types';

export const getTokenDecimals = ({
  pairs,
  identifier
}: {
  pairs: PairType[];
  identifier: string;
}) => {
  const isRewa =
    identifier === REWA_IDENTIFIER ||
    identifier.startsWith(`W${REWA_IDENTIFIER}-`);

  if (isRewa) {
    return DECIMALS;
  }

  let token: DcdtType | undefined;
  pairs.forEach(({ firstToken, secondToken }) => {
    if (firstToken.identifier === identifier) {
      token = firstToken;
      return;
    }

    if (secondToken.identifier === identifier) {
      token = secondToken;
      return;
    }
  });

  return token?.decimals ?? 0;
};
