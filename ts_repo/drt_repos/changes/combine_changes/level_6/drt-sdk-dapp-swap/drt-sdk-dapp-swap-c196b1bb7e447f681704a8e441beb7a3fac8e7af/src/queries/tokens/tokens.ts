import { gql } from '@apollo/client';
import {
  DcdtType,
  FactoryType,
  FilteredTokensType,
  UserDcdtType,
  WrappingInfoType
} from 'types';
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

export interface FilteredTokensQueryType {
  filteredTokens: FilteredTokensType;
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

export const GET_FILTERED_TOKENS = gql`
query swapPackageFilteredTokens ($enabledSwaps: Boolean, $pagination: ConnectionArgs, $searchInput: String, $identifiers: [String!]) {
    filteredTokens (pagination: $pagination, filters: {searchToken: $searchInput, enabledSwaps: $enabledSwaps, identifiers: $identifiers}) {
      edges {
        node {
          ${dcdtAttributes}
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

export const GET_FILTERED_TOKENS_AND_BALANCE = gql`
  query swapPackageFilteredTokensWithBalance ($identifiers: [String!], $pagination: ConnectionArgs, $searchInput: String, $offset: Int, $limit: Int, $enabledSwaps: Boolean) {
   filteredTokens (pagination: $pagination, filters: {searchToken: $searchInput, enabledSwaps: $enabledSwaps, identifiers: $identifiers}) {
      edges {
        node {
          ${dcdtAttributes}
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
