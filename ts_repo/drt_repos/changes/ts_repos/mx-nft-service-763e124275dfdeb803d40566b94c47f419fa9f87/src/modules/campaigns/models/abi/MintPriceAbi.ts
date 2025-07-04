import { BigUIntValue, TokenIdentifierValue } from '@terradharitri/sdk-core';

export interface MintPrice {
  token_id: TokenIdentifierValue;
  amount: BigUIntValue;
}
