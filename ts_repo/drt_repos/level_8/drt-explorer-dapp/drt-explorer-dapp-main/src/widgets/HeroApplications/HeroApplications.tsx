import { useSelector } from 'react-redux';

import { FormatREWA } from 'components';
import { useFetchGrowthEconomics } from 'hooks';
import { growthEconomicsSelector } from 'redux/selectors';
import { ChartContractsTransactions } from 'widgets';

export const HeroApplications = () => {
  const { applicationsDeployed, unprocessed } = useSelector(
    growthEconomicsSelector
  );

  const economics: any[] = [
    {
      label: 'Apps Deployed',
      value: <span className='text-primary'>{applicationsDeployed}</span>
    },
    {
      label: 'Fees Captured',
      value: (
        <FormatREWA value={unprocessed.feesCaptured} showLabel superSuffix />
      )
    },
    {
      label: 'Developer Rewards',
      value: (
        <FormatREWA
          value={unprocessed.developerRewards}
          showLabel
          superSuffix
        />
      )
    }
  ];

  useFetchGrowthEconomics();

  return (
    <div className='card-body p-0 '>
      <ChartContractsTransactions
        title='Transactions'
        showTransactions={false}
        showTotal={false}
        simpleTooltip={true}
        hasBorder={true}
        customStatistics={economics}
      />
    </div>
  );
};
