import { TokenArgumentType } from 'types/serverTransactions.types';

export interface DCDTValueDataType {
  tokenFormattedAmount: string | null;
  tokenExplorerLink: string;
  tokenLinkText: string;
  transactionTokens: TokenArgumentType[];
  token: TokenArgumentType;
  value: string | null;
  decimals: number | null;
  titleText: string;
}

export interface NFTValueDataType extends DCDTValueDataType {
  badgeText: string | null;
}

export interface TokenValueDataType extends DCDTValueDataType {
  showFormattedAmount: boolean;
}

export interface RewaValueDataType {
  value: string;
  formattedValue: string;
  decimals: number;
}
