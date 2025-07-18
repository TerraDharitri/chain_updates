import { numberToHex } from 'helpers/operations/toHex';
import { stringToHex } from 'helpers/operations/toHex';
import { prepareTransaction } from 'lib';
import { GAS_PRICE, SOVEREIGN_TRANSFER_GAS_LIMIT } from 'localConstants';
import {
  NftEnumType,
  PartialNftType,
  TokenType,
  DcdtEnumType,
  CollectionType
} from 'types';
import { RegisterTokenFormType } from '../types';

const TokenTypeMap: Record<string, number> = {
  [DcdtEnumType.FungibleDCDT]: 0,
  [NftEnumType.MetaDCDT]: 1,
  [NftEnumType.NonFungibleDCDT]: 2,
  [NftEnumType.SemiFungibleDCDT]: 3
};

export const getRegisterTokenTransaction = ({
  address,
  balance,
  nonce,
  values,
  token
}: {
  address: string;
  balance: string;
  nonce: number;
  values: RegisterTokenFormType;
  token: PartialNftType | TokenType | CollectionType;
}) => {
  const nft = token as PartialNftType;
  const isNft = Boolean(nft.nonce);
  const tokenIdentifier =
    'identifier' in token ? token.identifier : token.ticker;
  const tokenType = TokenTypeMap[nft.type] || 0;
  const tokenName = token.name;
  const tokenTicker = token.ticker?.split('-')[1];
  const tokenDecimals = token.decimals || 0;

  const args = [
    stringToHex(isNft ? nft.collection : tokenIdentifier),
    numberToHex(tokenType),
    stringToHex(tokenName),
    stringToHex(tokenTicker),
    numberToHex(tokenDecimals)
  ].join('@');

  const data = `registerToken@${args}`;

  return prepareTransaction({
    amount: '0.05',
    balance,
    chainId: values.chainId.value,
    data,
    gasLimit: SOVEREIGN_TRANSFER_GAS_LIMIT.toString(),
    gasPrice: GAS_PRICE.toString(),
    nonce,
    receiver: values.contract,
    sender: address
  });
};
