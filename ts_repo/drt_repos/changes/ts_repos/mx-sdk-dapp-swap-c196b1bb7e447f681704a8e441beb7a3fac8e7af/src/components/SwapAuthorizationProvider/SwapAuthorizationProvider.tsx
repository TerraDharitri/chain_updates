import React from 'react';
import { ApolloClient, from, HttpLink, InMemoryCache } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';
import { onError } from '@apollo/client/link/error';
import { print } from 'graphql';
import {
  SwapGraphQLAddressEnum,
  AuthorizationHeadersRequestParamsType
} from 'types';
import { AuthorizationContext } from './context';

export const SwapAuthorizationProvider = ({
  children,
  accessToken,
  graphQLAddress,
  isAccessTokenLoading = false,
  getAuthorizationHeaders
}: {
  accessToken?: string;
  children: React.ReactNode;
  isAccessTokenLoading?: boolean;
  graphQLAddress: SwapGraphQLAddressEnum | string;
  getAuthorizationHeaders?: (
    requestParams: AuthorizationHeadersRequestParamsType
  ) => Promise<void | null | Record<string, string>>;
}) => {
  const authMiddleware = setContext(async (req, { headers }) => {
    const requestParams: AuthorizationHeadersRequestParamsType = {
      url: graphQLAddress,
      params: req?.variables,
      body: {
        operationName: req?.operationName,
        variables: req?.variables,
        query: print(req?.query)
      },
      method: 'POST'
    };

    const authorizationHeaders = await getAuthorizationHeaders?.(requestParams);

    const authorizationBearerHeader = accessToken
      ? {
          Authorization: `Bearer ${accessToken}`
        }
      : {};

    return {
      headers: {
        ...headers,
        ...authorizationBearerHeader,
        ...authorizationHeaders
      }
    };
  });

  const errorLink = onError(({ graphQLErrors, networkError }) => {
    if (graphQLErrors) {
      graphQLErrors.forEach(({ message, locations, path }) => {
        console.log(
          `[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`
        );
      });
    }

    if (networkError) {
      console.log(`[Network error]: ${networkError}`);
    }
  });

  const httpLink = from([errorLink, new HttpLink({ uri: graphQLAddress })]);

  const client = new ApolloClient({
    cache: new InMemoryCache(),
    link: authMiddleware.concat(httpLink),
    queryDeduplication: false, // FIX: fixes canceled queries not beeing sent to the server when retriggered
    defaultOptions: {
      watchQuery: {
        fetchPolicy: 'no-cache',
        errorPolicy: 'all'
      },
      query: {
        fetchPolicy: 'no-cache',
        errorPolicy: 'all'
      }
    }
  });

  return (
    <AuthorizationContext.Provider
      value={{
        client,
        accessToken,
        isAccessTokenLoading,
        isAuthenticated: Boolean(accessToken)
      }}
    >
      {children}
    </AuthorizationContext.Provider>
  );
};
