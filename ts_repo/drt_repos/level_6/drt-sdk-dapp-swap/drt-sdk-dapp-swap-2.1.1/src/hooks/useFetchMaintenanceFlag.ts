import { useMemo } from 'react';
import { useAuthorizationContext } from 'components/SwapAuthorizationProvider';
import { useQueryWrapper } from 'hooks';
import { maintenanceQuery, MaintenanceQueryType } from 'queries';

export const useFetchMaintenanceFlag = (isPollingEnabled?: boolean) => {
  const { client } = useAuthorizationContext();

  const { data } = useQueryWrapper<MaintenanceQueryType>({
    isPollingEnabled,
    queryOptions: { variables: {}, client },
    query: maintenanceQuery
  });

  const isMaintenance = useMemo(
    () => data?.factory.maintenance,
    [data?.factory.maintenance]
  );

  return isMaintenance;
};
