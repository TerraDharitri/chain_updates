import { useSelector } from 'react-redux';

import { Tabs } from 'components/Tabs';
import { urlBuilder } from 'helpers';
import { activeNetworkSelector } from 'redux/selectors';

export const NativeTokenTabs = () => {
  const { rewaLabel = 'rewa' } = useSelector(activeNetworkSelector);

  const tabs = [
    {
      tabLabel: 'Transactions',
      tabTo: urlBuilder.nativeTokenDetails(rewaLabel),
      activationRoutes: [urlBuilder.nativeTokenDetails(rewaLabel)]
    },
    {
      tabLabel: 'Holders',
      tabTo: urlBuilder.nativeTokenDetailsAccounts(rewaLabel),
      activationRoutes: [urlBuilder.nativeTokenDetailsAccounts(rewaLabel)]
    }
  ];

  return <Tabs tabs={tabs} />;
};
