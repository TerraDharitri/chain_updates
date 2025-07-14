import { DECIMALS, formatAmount } from 'lib/sdkDappUtils';

export const getRewaValueData = (value: string) => ({
  rewaValueData: {
    value,
    formattedValue: formatAmount({ input: value }),
    decimals: DECIMALS
  }
});
