import { Component, h, Prop } from '@stencil/core';
import classNames from 'classnames';
import { DecodeMethodEnum } from 'components/functional/sign-transactions-panel/sign-transactions-panel.types';

const signTransactionsAdvancedDataDecodeClasses: Record<string, string> = {
  icon: 'drt:transition-all drt:duration-200 drt:ease-in-out drt:relative drt:h-3! drt:w-auto!',
  iconRotated: 'drt:rotate-90',
};

@Component({
  tag: 'drt-sign-transactions-advanced-data-decode',
  styleUrl: 'sign-transactions-advanced-data-decode.scss',
  shadow: true,
})
export class SignTransactionsAdvancedDataDecode {
  @Prop() isToggled: boolean = false;
  @Prop() currentDecodeMethod: DecodeMethodEnum = DecodeMethodEnum.decimal;

  render() {
    return (
      <div class="sign-transactions-advanced-data-decode">
        <div class="sign-transactions-advanced-data-decode-label">{this.currentDecodeMethod}</div>

        <drt-single-angle-down-icon
          class={classNames('sign-transactions-advanced-data-decode-icon', {
            [signTransactionsAdvancedDataDecodeClasses.icon]: true,
            [signTransactionsAdvancedDataDecodeClasses.iconRotated]: this.isToggled,
          })}
        />
      </div>
    );
  }
}
