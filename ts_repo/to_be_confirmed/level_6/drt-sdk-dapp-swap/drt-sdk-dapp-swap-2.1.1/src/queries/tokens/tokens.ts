import { gql } from '@apollo/client';
import { DcdtType, FactoryType, UserDcdtType, WrappingInfoType } from 'types';
import {
  dcdtAttributes,
  factoryAttributes,
  userDcdtAttributes
} from '../attributes';

export interface TokensType {
  tokens: DcdtType[];
  factory: FactoryType;
  userTokens?: UserDcdtType[];
  wrappingInfo: WrappingInfoType[];
}

export const GET_TOKENS = gql`
  query swapPackageTokens ($identifiers: [String!], $enabledSwaps: Boolean) {
    tokens(identifiers: $identifiers, enabledSwaps: $enabledSwaps) {
      ${dcdtAttributes}
    }
    wrappingInfo {
      wrappedToken {
        ${dcdtAttributes}
      }
    }
    factory {
      ${factoryAttributes}
    }
  }
`;

export const GET_TOKENS_AND_BALANCE = gql`
  query swapPackageTokensWithBalance ($identifiers: [String!], $offset: Int, $limit: Int, $enabledSwaps: Boolean) {
    tokens(identifiers: $identifiers, enabledSwaps: $enabledSwaps) {
      ${dcdtAttributes}
    }
    userTokens (offset: $offset, limit: $limit) {
      ${userDcdtAttributes}
    }
    wrappingInfo {
      wrappedToken {
        ${dcdtAttributes}
      }
    }
    factory {
      ${factoryAttributes}
    }
  }
`;
