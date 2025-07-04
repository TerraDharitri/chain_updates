import classNames from 'classnames';

import { getNftText } from 'helpers';
import { NftTypeEnum, TokenTypeEnum } from 'types';

export const NftTypeBadge = ({
  type,
  className
}: {
  type: NftTypeEnum | TokenTypeEnum;
  className?: string;
}) => {
  switch (type) {
    // default NFT types
    case NftTypeEnum.SemiFungibleDCDT:
    case NftTypeEnum.NonFungibleDCDT:
    case NftTypeEnum.MetaDCDT:
      return (
        <div
          className={classNames(
            'badge badge-outline',
            { 'badge-outline-orange': type === NftTypeEnum.SemiFungibleDCDT },
            { 'badge-outline-yellow': type === NftTypeEnum.NonFungibleDCDT },
            { 'badge-outline-green': type === NftTypeEnum.MetaDCDT },
            className
          )}
        >
          {getNftText(type)}
        </div>
      );

    default:
      return null;
  }
};
