import React, { useEffect } from 'react';
import { useFormik } from 'formik';
import {
  SwapFormValuesEnum,
  UserDcdtType,
  FactoryType,
  SelectOptionType,
  SwapRouteType
} from 'types';
import { meaningfulFormatAmount } from 'utils';
import { useSwapValidationSchema } from 'validation/hooks/useSwapValidationSchema';
import { TokenSelect } from '../TokenSelect/TokenSelect';
import { SwapFormInputValidationErrorDisplay } from './SwapFormInputValidationErrorDisplay';

export interface SwapFormType {
  tokens: UserDcdtType[];
  firstAmount: string;
  secondAmount: string;
  firstToken?: SelectOptionType;
  secondToken?: SelectOptionType;
  activeRoute?: SwapRouteType;
  swapConfig?: FactoryType;
  handleOnSubmit?: () => void;
  handleSwitchTokens: () => void;
  handleOnFirstMaxBtnChange: () => void;
  handleOnChangeFirstAmount: (amount: string) => void;
  handleOnChangeSecondAmount: (amount: string) => void;
  handleOnChangeFirstSelect: (option?: SelectOptionType) => void;
  handleOnChangeSecondSelect: (option?: SelectOptionType) => void;
}

export const SwapForm = ({
  tokens,
  swapConfig,
  firstToken,
  firstAmount,
  secondToken,
  secondAmount,
  activeRoute,
  handleOnChangeFirstAmount,
  handleOnChangeSecondAmount,
  handleOnChangeFirstSelect,
  handleOnChangeSecondSelect,
  handleOnFirstMaxBtnChange,
  handleSwitchTokens,
  handleOnSubmit
}: SwapFormType) => {
  const initialValues = {
    firstAmount: '',
    firstToken: firstToken,
    secondAmount: '',
    secondToken: secondToken,
    activeRoute: activeRoute
  };

  const validationSchema = useSwapValidationSchema({
    firstToken,
    secondToken,
    minAcceptedAmount: swapConfig?.minSwapAmount
  });

  const onSubmit = () => {
    if (handleOnSubmit) {
      handleOnSubmit();
    }
  };

  const {
    handleSubmit,
    handleBlur,
    handleChange,
    setFieldValue,
    setTouched,
    errors,
    touched
  } = useFormik<typeof initialValues>({
    onSubmit,
    initialValues,
    validationSchema
  });

  const selectOptions: SelectOptionType[] = tokens?.map((token) => {
    return { label: token.identifier, value: token.identifier, token };
  });

  const resetSwapForm = () => {
    setTouched({}, false);
    setFieldValue(SwapFormValuesEnum.firstAmount, '');
    setFieldValue(SwapFormValuesEnum.secondAmount, '');
    setFieldValue(SwapFormValuesEnum.activeRoute, undefined);
  };

  const onChangeFirstSelect = (option?: SelectOptionType) => {
    setFieldValue(SwapFormValuesEnum.firstToken, option);
    resetSwapForm();

    handleOnChangeFirstSelect(option);
  };

  const onChangeSecondSelect = (option?: SelectOptionType) => {
    setFieldValue(SwapFormValuesEnum.secondToken, option);
    resetSwapForm();

    handleOnChangeSecondSelect(option);
  };

  const firstAmountOnChange = (event: React.FormEvent<HTMLInputElement>) => {
    handleChange(event);
    handleOnChangeFirstAmount(event.currentTarget.value);
  };

  const secondAmountOnChange = (event: React.FormEvent<HTMLInputElement>) => {
    handleChange(event);
    handleOnChangeSecondAmount(event.currentTarget.value);
  };

  const firstBalance = firstToken?.token
    ? meaningfulFormatAmount({
        amount: firstToken.token.balance ?? '0',
        decimals: firstToken?.token.decimals
      })
    : 0;

  const secondBalance = secondToken?.token
    ? meaningfulFormatAmount({
        amount: secondToken.token.balance ?? '0',
        decimals: secondToken?.token.decimals
      })
    : 0;

  useEffect(() => {
    setFieldValue(SwapFormValuesEnum.firstAmount, firstAmount, true);
  }, [firstAmount]);

  useEffect(() => {
    setFieldValue(SwapFormValuesEnum.secondAmount, secondAmount, true);
  }, [secondAmount]);

  useEffect(() => {
    setFieldValue(SwapFormValuesEnum.activeRoute, activeRoute, true);
  }, [activeRoute]);

  return (
    <form
      className='dapp-core-swap-form'
      noValidate={true}
      onSubmit={handleSubmit}
    >
      <div>
        <TokenSelect
          id={SwapFormValuesEnum.firstToken}
          name={SwapFormValuesEnum.firstToken}
          value={firstToken}
          isSearchable={false}
          options={selectOptions}
          onChange={onChangeFirstSelect}
          onBlur={handleBlur}
          disabledOption={secondToken}
        />

        <div>Balance: {firstBalance}</div>

        <input
          type='number'
          min={0}
          step='0.000001'
          name={SwapFormValuesEnum.firstAmount}
          value={firstAmount}
          onChange={firstAmountOnChange}
          onBlur={handleBlur}
        />

        <SwapFormInputValidationErrorDisplay
          fieldName={SwapFormValuesEnum.firstAmount}
          errors={errors}
          touched={touched}
        />

        <button type='button' onClick={handleOnFirstMaxBtnChange}>
          max
        </button>
      </div>

      <br />
      <button type='button' onClick={handleSwitchTokens}>
        {' '}
        {'> reverse <'}{' '}
      </button>
      <br />
      <br />

      <div>
        <TokenSelect
          id={SwapFormValuesEnum.secondToken}
          name={SwapFormValuesEnum.secondToken}
          value={secondToken}
          isSearchable={false}
          options={selectOptions}
          onChange={onChangeSecondSelect}
          disabledOption={firstToken}
        />

        <div>Balance: {secondBalance}</div>

        <input
          type='number'
          min={0}
          step='0.000001'
          name={SwapFormValuesEnum.secondAmount}
          value={secondAmount}
          onChange={secondAmountOnChange}
          onBlur={handleBlur}
        />

        <SwapFormInputValidationErrorDisplay
          fieldName={SwapFormValuesEnum.secondAmount}
          errors={errors}
          touched={touched}
        />
      </div>

      <br />

      <button type='submit'>Swap</button>
    </form>
  );
};
