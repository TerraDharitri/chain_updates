import { gql } from '@apollo/client';
import {
  EsdtType,
  FactoryType,
  FilteredTokensType,
  UserEsdtType,
  WrappingInfoType
} from 'types';
import {
  esdtAttributes,
  factoryAttributes,
  userEsdtAttributes
} from '../attributes';

export interface TokensType {
  tokens: EsdtType[];
  factory: FactoryType;
  userTokens?: UserEsdtType[];
  wrappingInfo: WrappingInfoType[];
}

export interface FilteredTokensQueryType {
  filteredTokens: FilteredTokensType;
  factory: FactoryType;
  userTokens?: UserEsdtType[];
  wrappingInfo: WrappingInfoType[];
}

export const GET_TOKENS = gql`
  query swapPackageTokens ($identifiers: [String!], $enabledSwaps: Boolean) {
    tokens(identifiers: $identifiers, enabledSwaps: $enabledSwaps) {
      ${esdtAttributes}
    }
    wrappingInfo {
      wrappedToken {
        ${esdtAttributes}
      }
    }
    factory {
      ${factoryAttributes}
    }
  }
`;

export const GET_FILTERED_TOKENS = gql`
query swapPackageFilteredTokens ($enabledSwaps: Boolean, $pagination: ConnectionArgs, $searchInput: String, $identifiers: [String!]) {
    filteredTokens (pagination: $pagination, filters: {searchToken: $searchInput, enabledSwaps: $enabledSwaps, identifiers: $identifiers}) {
      edges {
        node {
          ${esdtAttributes}
        }
        cursor
      }
      pageInfo {
        hasNextPage
      }
      pageData {
        count
      }
    }
    wrappingInfo {
      wrappedToken {
        ${esdtAttributes}
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
      ${esdtAttributes}
    }
    userTokens (offset: $offset, limit: $limit) {
      ${userEsdtAttributes}
    }
    wrappingInfo {
      wrappedToken {
        ${esdtAttributes}
      }
    }
    factory {
      ${factoryAttributes}
    }
  }
`;

export const GET_FILTERED_TOKENS_AND_BALANCE = gql`
  query swapPackageFilteredTokensWithBalance ($identifiers: [String!], $pagination: ConnectionArgs, $searchInput: String, $offset: Int, $limit: Int, $enabledSwaps: Boolean) {
   filteredTokens (pagination: $pagination, filters: {searchToken: $searchInput, enabledSwaps: $enabledSwaps, identifiers: $identifiers}) {
      edges {
        node {
          ${esdtAttributes}
        }
        cursor
      }
      pageInfo {
        hasNextPage
      }
      pageData {
        count
      }
    }
    userTokens (offset: $offset, limit: $limit) {
      ${userEsdtAttributes}
    }
    wrappingInfo {
      wrappedToken {
        ${esdtAttributes}
      }
    }
    factory {
      ${factoryAttributes}
    }
  }
`;
