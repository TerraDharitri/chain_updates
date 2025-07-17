import { useSelector } from 'react-redux';

import { Tabs } from 'components/Tabs';
import { urlBuilder, getNftText } from 'helpers';
import { collectionSelector } from 'redux/selectors';
import { collectionRoutes, tokensRoutes } from 'routes';
import { NftTypeEnum } from 'types';

export const CollectionTabs = () => {
  const { collectionState } = useSelector(collectionSelector);
  const { collection, roles, type } = collectionState;

  const isMetaDCDT = type && type === NftTypeEnum.MetaDCDT;

  const tabs = [
    {
      show: !isMetaDCDT,
      tabTo: urlBuilder.collectionDetails(collection),
      tabLabel: `${getNftText(type)}s`,
      activationRoutes: [collectionRoutes.collectionDetails]
    },
    {
      tabTo: isMetaDCDT
        ? urlBuilder.tokenMetaDcdtDetails(collection)
        : urlBuilder.collectionDetailsTransactions(collection),
      tabLabel: 'Transactions',
      activationRoutes: [
        collectionRoutes.collectionDetailsTransactions,
        tokensRoutes.tokensMetaDcdtDetails
      ]
    },
    {
      show: Boolean(roles),
      tabTo: isMetaDCDT
        ? urlBuilder.tokenMetaDcdtDetailsRoles(collection)
        : urlBuilder.collectionDetailsRoles(collection),
      tabLabel: 'Roles',
      activationRoutes: [
        collectionRoutes.collectionDetailsRoles,
        tokensRoutes.tokensMetaDcdtDetailsRoles
      ]
    }
  ];

  return <Tabs tabs={tabs} />;
};
