import { gql } from '@apollo/client';
import { transactionAttributes } from 'queries/attributes/transaction';

export const wrapEgldQuery = gql`
  query swapPackageWrapEgld ($wrappingAmount: String!) {
    wrapEgld(amount: $wrappingAmount) {
      ${transactionAttributes}
    }
  }
`;
