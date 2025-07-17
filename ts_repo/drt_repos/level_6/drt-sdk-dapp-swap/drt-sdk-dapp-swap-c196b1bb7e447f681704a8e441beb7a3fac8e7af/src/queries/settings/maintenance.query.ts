import { gql } from '@apollo/client';

export interface MaintenanceQueryType {
  factory: {
    maintenance: boolean;
  };
}

export const maintenanceQuery = gql`
  query swapPackageMaintenance {
    factory {
      maintenance
    }
  }
`;
