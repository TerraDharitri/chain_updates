import { useEffect } from 'react';
import uniqBy from 'lodash/unionBy';
import { useGetAccountInfo, useGetNetworkConfig } from 'lib';
import { useLazyGetTokensQuery } from 'redux/endpoints';
import { TokenType } from 'types';

const defaultValues = {
  owner: '',
  minted: '',
  burnt: '',
  supply: '',
  circulatingSupply: '',
  canBurn: false,
  canChangeOwner: false,
  canFreeze: false,
  canMint: false,
  canPause: false,
  canUpgrade: false,
  canWipe: false,
  isPaused: false,
  transactions: 0,
  accounts: 0
};

export const useGetTokensWithRewa = () => {
  const { websocketEvent, address, account } = useGetAccountInfo();
  const { network } = useGetNetworkConfig();
  const [fetchTokens, { data: tokens, isLoading }] = useLazyGetTokensQuery();

  const rewaToken: TokenType = {
    ...defaultValues,
    identifier: network.rewaLabel,
    name: network.rewaLabel,
    balance: account?.balance
  };

  useEffect(() => {
    fetchTokens(address);
  }, [address, websocketEvent]);

  const usedTokens = [rewaToken, ...(tokens ?? [])];

  return {
    tokens: uniqBy(usedTokens, 'identifier'),
    isLoading
  };
};
