import React from 'react';
import { DcdtType } from 'types';
import { meaningfulFormatAmount, roundAmount } from 'utils';

interface SelectOptionType {
  label: string;
  value: string;
  token: DcdtType;
  isDisabled: boolean;
  inDropdown?: boolean;
  handleDisabledOptionClick?: any;
}

export const SelectOption = ({
  value,
  token,
  isDisabled,
  inDropdown = false,
  handleDisabledOptionClick
}: SelectOptionType) => {
  const handleOnClick = (e: any) => {
    return;
    if (isDisabled && handleDisabledOptionClick) {
      handleDisabledOptionClick(e);
    }
  };

  return (
    <div
      className={`dapp-core-swap-select-option ${isDisabled ? 'disabled' : ''}`}
      onClick={handleOnClick}
    >
      <div className='d-flex flex-row align-items-center'>
        <div
          className={`token-image mr-2 ${
            inDropdown ? 'in-dropdown' : 'd-none d-md-flex'
          }
        `}
        >
          <img
            src={token?.assets?.svgUrl}
            alt={value}
            className='token-symbol'
            style={{ width: '1.5rem', height: '1.5rem' }}
          />
        </div>

        <div className='d-flex flex-column'>
          {token.ticker}{' '}
          <small className='text-secondary'>
            {roundAmount(token.price ?? '0')}
          </small>
        </div>
      </div>

      {inDropdown && (
        <div className='d-flex flex-column ml-spacer ml-lg-5 align-items-end'>
          {meaningfulFormatAmount({
            amount: token.balance ?? '',
            decimals: token.decimals
          })}

          {/* {token.totalUsdPrice && (
            <small className='text-secondary'>
              {token.totalUsdPrice !== '$0' ? <>≈&nbsp;</> : <></>}
              {token.totalUsdPrice}
            </small>
          )} */}
        </div>
      )}
    </div>
  );
};
