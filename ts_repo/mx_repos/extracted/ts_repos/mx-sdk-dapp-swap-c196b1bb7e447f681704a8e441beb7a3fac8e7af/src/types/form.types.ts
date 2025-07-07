import { AssetsType, UserEsdtType } from './tokens.types';

export interface SelectOptionType {
  label: string;
  value: string;
  token: UserEsdtType;
  assets?: AssetsType;
}

export type TokenOptionType = SelectOptionType;

export enum SwapFormValuesEnum {
  firstToken = 'firstToken',
  secondToken = 'secondToken',
  firstAmount = 'firstAmount',
  secondAmount = 'secondAmount',
  activeRoute = 'activeRoute'
}
