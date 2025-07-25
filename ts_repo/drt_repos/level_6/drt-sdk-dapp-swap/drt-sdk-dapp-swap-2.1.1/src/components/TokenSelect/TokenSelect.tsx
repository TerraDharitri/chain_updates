import React from 'react';
import Select from 'react-select';
import { SelectOptionType } from 'types';
import { SelectOption } from '../SelectOption';

interface TokenSelectType {
  id?: string;
  value?: any;
  name: string;
  options: any;
  defaultValue?: any;
  className?: string;
  disabled?: boolean;
  isSearchable?: boolean;
  noOptionsMessage?: string;
  handleDisabledOptionClick?: any;
  disabledOption?: SelectOptionType;
  onChange: (option: any) => void;
  onBlur?: (option: any) => void;
  onFocus?: (props: any) => void;
}

export const TokenSelect = ({
  id,
  name,
  value,
  options,
  defaultValue,
  isSearchable,
  className = '',
  disabledOption,
  disabled = false,
  handleDisabledOptionClick,
  noOptionsMessage = 'No Tokens',
  onChange,
  onBlur,
  onFocus
}: TokenSelectType) => {
  const ref = React.useRef(null);
  const disableOption = (option: any) => option.value === disabledOption?.value;

  const FormatOptionLabel =
    () =>
    (option: any, { context }: { context: any }) => {
      const { label, value: val, token } = option;
      const inDropdown = context === 'menu';
      const isDisabled = disableOption(option);

      const args = {
        inDropdown,
        label,
        value: val,
        token,
        isDisabled,
        handleDisabledOptionClick
      };

      return <SelectOption {...args} />;
    };

  return (
    <div className='dapp-core-swap-select-container'>
      <Select
        ref={ref}
        placeholder='Select token'
        id={id}
        name={name}
        value={value}
        options={options}
        className={`dapp-core-swap-select ${className}`}
        isOptionDisabled={disableOption}
        isDisabled={disabled}
        noOptionsMessage={() => noOptionsMessage || 'No options'}
        maxMenuHeight={260}
        onChange={(e) => {
          onChange(e);
          if (ref && ref.current !== null) {
            (ref.current as any).blur();
          }
        }}
        onBlur={onBlur}
        onFocus={onFocus}
        isSearchable={isSearchable}
        defaultValue={defaultValue}
        formatOptionLabel={FormatOptionLabel()}
      />
    </div>
  );
};
