import { useEffect, useMemo } from 'react';
import { useGetTokensWithRewa } from 'hooks';
import { getRewaLabel, useGetAccountInfo } from 'lib';
import { useLazyGetNftsQuery } from 'redux/endpoints';
import { SendTypeEnum, TokenOptionType } from 'types';

export const useTokenOptions = ({
  sendType,
  skipAddRewa
}: {
  sendType: SendTypeEnum;
  skipAddRewa?: boolean;
}) => {
  const { address, websocketEvent } = useGetAccountInfo();
  const { tokens, isLoading: isLoadingTokens } = useGetTokensWithRewa();
  const [fetchNFTs, { data: nfts, isLoading: isLoadingNfts }] =
    useLazyGetNftsQuery();
  const rewaLabel = getRewaLabel();

  const getTokenOptionsByType = (type: SendTypeEnum): TokenOptionType[] => {
    if (type === SendTypeEnum.nft) {
      return (
        nfts?.map((token) => ({
          value: token.identifier,
          label: token.name
        })) ?? []
      );
    }

    return tokens
      .filter((token) => !skipAddRewa || token.identifier !== rewaLabel)
      .map((token) => ({
        value: token.identifier,
        label: token.name
      }));
  };

  const getTokens = (type: SendTypeEnum) =>
    type === SendTypeEnum.nft ? nfts : tokens;

  useEffect(() => {
    fetchNFTs({ address });
  }, [address, websocketEvent]);

  const tokenOptions = useMemo(
    () => getTokenOptionsByType(sendType),
    [nfts, tokens, sendType]
  );

  const allTokens = [...tokens, ...(nfts || [])].filter(
    (token) => !skipAddRewa || token.identifier !== rewaLabel
  );

  return {
    allTokens,
    getTokenOptionsByType,
    getTokens,
    isLoading: isLoadingTokens || isLoadingNfts,
    tokenOptions,
    tokens: getTokens(sendType) || []
  };
};
