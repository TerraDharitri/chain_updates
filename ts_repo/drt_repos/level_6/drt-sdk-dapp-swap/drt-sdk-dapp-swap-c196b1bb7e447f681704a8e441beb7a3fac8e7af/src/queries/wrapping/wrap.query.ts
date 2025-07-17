import { gql } from '@apollo/client';
import { transactionAttributes } from 'queries/attributes/transaction';

export const wrapRewaQuery = gql`
  query swapPackageWrapRewa ($wrappingAmount: String!) {
    wrapRewa(amount: $wrappingAmount) {
      ${transactionAttributes}
    }
  }
`;
