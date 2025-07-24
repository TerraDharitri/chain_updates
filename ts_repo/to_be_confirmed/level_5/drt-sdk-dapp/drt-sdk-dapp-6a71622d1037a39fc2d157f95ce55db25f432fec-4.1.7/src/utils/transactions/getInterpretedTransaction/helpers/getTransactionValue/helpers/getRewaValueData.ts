import { DECIMALS } from 'constants/index';
import { formatAmount } from 'utils/operations/formatAmount';

export const getRewaValueData = (value: string) => ({
  rewaValueData: {
    value,
    formattedValue: formatAmount({ input: value }),
    decimals: DECIMALS
  }
});
