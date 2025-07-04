import { AssetsType, UserDcdtType } from './tokens.types';

export interface SelectOptionType {
  label: string;
  value: string;
  token: UserDcdtType;
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
