import { WithClassnameType } from '../types/withClassname.types';

export interface FormatAmountPropsType extends WithClassnameType {
  value: string;
  showLastNonZeroDecimal?: boolean;
  showLabel?: boolean;
  token?: string;
  digits?: number;
  decimals?: number;
  rewaLabel?: string;
  'data-testid'?: string;
}
