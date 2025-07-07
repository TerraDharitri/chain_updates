import { EGLD_IDENTIFIER } from 'constants/general';
import { EsdtType, SwapActionTypesEnum } from 'types';

export const getSwapActionType = ({
  firstTokenId,
  secondTokenId,
  wrappedEgld
}: {
  firstTokenId?: string;
  secondTokenId?: string;
  wrappedEgld?: EsdtType;
}) => {
  if (!firstTokenId || !secondTokenId || !wrappedEgld) {
    return;
  }

  const isWrapEgldTransaction =
    firstTokenId === EGLD_IDENTIFIER &&
    secondTokenId === wrappedEgld?.identifier;

  const isUnwrapEgldTransaction =
    firstTokenId === wrappedEgld?.identifier &&
    secondTokenId === EGLD_IDENTIFIER;

  if (isWrapEgldTransaction) return SwapActionTypesEnum.wrap;
  if (isUnwrapEgldTransaction) return SwapActionTypesEnum.unwrap;

  return SwapActionTypesEnum.swap;
};
