import { REWA_IDENTIFIER } from 'constants/general';
import { DcdtType, SwapActionTypesEnum } from 'types';

export const getSwapActionType = ({
  firstTokenId,
  secondTokenId,
  wrappedRewa
}: {
  firstTokenId?: string;
  secondTokenId?: string;
  wrappedRewa?: DcdtType;
}) => {
  if (!firstTokenId || !secondTokenId || !wrappedRewa) {
    return;
  }

  const isWrapRewaTransaction =
    firstTokenId === REWA_IDENTIFIER &&
    secondTokenId === wrappedRewa?.identifier;

  const isUnwrapRewaTransaction =
    firstTokenId === wrappedRewa?.identifier &&
    secondTokenId === REWA_IDENTIFIER;

  if (isWrapRewaTransaction) return SwapActionTypesEnum.wrap;
  if (isUnwrapRewaTransaction) return SwapActionTypesEnum.unwrap;

  return SwapActionTypesEnum.swap;
};
