import React from 'react';
import { MAINNET_REWA_LABEL } from 'constants/index';
import { FormatAmount } from 'UI/FormatAmount/FormatAmount';
import { FormatAmountPropsType } from 'UI/FormatAmount/formatAmount.types';

let deprecationMessageDisplayed = false;

/**
 * !!! This function is deprecated. Please use formatAmount instead.
 * @param props.rewaLabel  if not provided, will fallback on **REWA**
 */
export const Denominate = (props: FormatAmountPropsType) => {
  if (!deprecationMessageDisplayed) {
    console.warn(
      '!!! Be aware !!! The "Denominate" component is deprecated. Please use "FormatAmount" instead.'
    );

    deprecationMessageDisplayed = true;
  }

  const rewaLabel = props.rewaLabel || MAINNET_REWA_LABEL;

  const formatAmountProps = { ...props, rewaLabel };

  return <FormatAmount {...formatAmountProps} />;
};
