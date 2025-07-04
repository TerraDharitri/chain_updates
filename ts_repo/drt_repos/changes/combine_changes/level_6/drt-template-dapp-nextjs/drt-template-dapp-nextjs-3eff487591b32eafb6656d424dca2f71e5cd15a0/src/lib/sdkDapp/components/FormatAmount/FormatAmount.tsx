import { WithClassnameType } from '@/types';
import {
  MvxFormatAmount,
  MvxFormatAmountPropsType,
  DECIMALS,
  DIGITS,
  FormatAmountController,
  useGetNetworkConfig
} from '@/lib';

interface FormatAmountPropsType
  extends Partial<MvxFormatAmountPropsType>,
    WithClassnameType {
  rewaLabel?: string;
  value: string;
}

export const FormatAmount = (props: FormatAmountPropsType) => {
  const {
    network: { rewaLabel }
  } = useGetNetworkConfig();

  const { isValid, valueDecimal, valueInteger, label } =
    FormatAmountController.getData({
      digits: DIGITS,
      decimals: DECIMALS,
      rewaLabel,
      ...props,
      input: props.value
    });

  return (
    <MvxFormatAmount
      class={props.className}
      dataTestId={props['data-testid']}
      isValid={isValid}
      label={label}
      showLabel={props.showLabel}
      valueDecimal={valueDecimal}
      valueInteger={valueInteger}
    />
  );
};
