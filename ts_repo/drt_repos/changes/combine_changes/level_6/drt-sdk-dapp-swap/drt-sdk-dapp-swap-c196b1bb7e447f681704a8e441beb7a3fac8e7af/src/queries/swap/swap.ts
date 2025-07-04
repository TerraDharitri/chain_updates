import { gql } from '@apollo/client';
import { GraphqlErrorsResponseType, SwapRouteType } from 'types';
import { transactionAttributes, pairAttributes } from '../attributes';

const transactionsPlaceholder = 'TRANSACTIONS_PLACEHOLDER';

const transactionsString = `
  transactions {
    ${transactionAttributes}
  }
`;

const swapString = `
  query swapPackageSwapRoute (
    $amountIn: String
    $amountOut: String
    $tokenInID: String!
    $tokenOutID: String!
    $tolerance: Float!
  ) {
    swap(
      amountIn: $amountIn
      amountOut: $amountOut
      tokenInID: $tokenInID
      tokenOutID: $tokenOutID
      tolerance: $tolerance
    ) {
      amountIn
      tokenInID
      tokenInPriceUSD
      tokenInExchangeRateDenom

      amountOut
      tokenOutID
      tokenOutPriceUSD
      tokenOutExchangeRateDenom

      fees
      swapType
      tokenRoute
      pricesImpact

      maxPriceDeviationPercent
      tokensPriceDeviationPercent

      intermediaryAmounts
      pairs {
        ${pairAttributes}
      }
      ${transactionsPlaceholder}
    }
  }
`;

export const swapQuery = gql`
  ${swapString.replace(transactionsPlaceholder, transactionsString)}
`;

export const swapWithoutTransactionsQuery = gql`
  ${swapString.replace(transactionsPlaceholder, '')}
`;

export interface SwapRouteQueryResponseType {
  swap: SwapRouteType;
}
