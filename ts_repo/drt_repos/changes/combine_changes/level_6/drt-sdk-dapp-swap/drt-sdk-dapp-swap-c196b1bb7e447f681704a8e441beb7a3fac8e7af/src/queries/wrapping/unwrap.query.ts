import { gql } from '@apollo/client';
import { transactionAttributes } from 'queries/attributes/transaction';

export const unwrapRewaQuery = gql`
  query swapPackageUnwrapRewa ($wrappingAmount: String!) {
    unwrapRewa(amount: $wrappingAmount) {
      ${transactionAttributes}
    }
  }
`;
