import classNames from 'classnames';

import { NftSubtypeEnum } from 'types';

export const NftSubTypeBadge = ({
  subType,
  className
}: {
  subType: NftSubtypeEnum;
  className?: string;
}) => {
  switch (subType) {
    // NFT Subtypes
    case NftSubtypeEnum.DynamicSemiFungibleDCDT:
    case NftSubtypeEnum.DynamicNonFungibleDCDT:
    case NftSubtypeEnum.NonFungibleDCDTv2:
    case NftSubtypeEnum.DynamicMetaDCDT:
      return (
        <div
          className={classNames(
            'badge',
            {
              'badge-orange text-orange-100':
                subType === NftSubtypeEnum.DynamicSemiFungibleDCDT,
              'badge-yellow text-orange-100':
                subType === NftSubtypeEnum.DynamicNonFungibleDCDT ||
                subType === NftSubtypeEnum.NonFungibleDCDTv2,
              'badge-green text-green-100':
                subType === NftSubtypeEnum.DynamicMetaDCDT
            },
            className
          )}
        >
          {subType}
        </div>
      );

    default:
      return null;
  }
};
