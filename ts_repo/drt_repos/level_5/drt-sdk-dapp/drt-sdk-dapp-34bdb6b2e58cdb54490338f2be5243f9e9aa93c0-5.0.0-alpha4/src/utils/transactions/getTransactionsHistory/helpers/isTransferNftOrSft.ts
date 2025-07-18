import { NftEnumType } from 'types/tokens.types';

export const isTransferNftOrSft = (
  transactionTransfer: Record<string, string>
) => {
  const assetTypesWithoutDenomination = [
    NftEnumType.SemiFungibleDCDT,
    NftEnumType.NonFungibleDCDT
  ];

  const transferType = transactionTransfer.type as NftEnumType;
  const isTransferTypeNftOrSFt =
    assetTypesWithoutDenomination.includes(transferType);

  return isTransferTypeNftOrSFt;
};
