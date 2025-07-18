import type { EventEmitter } from '@stencil/core';
import { Component, Event, h, Prop } from '@stencil/core';
import classNames from 'classnames';
import { DataTestIdsEnum } from 'constants/dataTestIds.enum';

import type { IAccountScreenData, ILedgerAccount } from '../../ledger-connect.types';

const TOTAL_ADDRESSES_COUNT = 5000;

@Component({
  tag: 'drt-ledger-addresses',
  styleUrl: 'ledger-addresses.scss',
  shadow: true,
})
export class LedgerAddresses {
  @Prop() accountScreenData: IAccountScreenData;
  @Prop() selectedIndex: number;

  @Event() accessWallet: EventEmitter;
  @Event({ bubbles: false, composed: false }) selectAccount: EventEmitter;
  @Event() pageChange: EventEmitter<number>;

  handleAccessWallet(event: MouseEvent) {
    event.preventDefault();
    this.accessWallet.emit();
  }

  handleSelectAccount(accountDerivationIndex: number) {
    return (event: MouseEvent) => {
      event.preventDefault();
      this.selectAccount.emit(accountDerivationIndex);
    };
  }

  private handlePageChange(event: CustomEvent) {
    event.preventDefault();
    event.stopPropagation();
    this.pageChange.emit(event.detail);
  }

  private processLedgerAddressIndex(accountDerivation: ILedgerAccount) {
    return Number(accountDerivation.index + 1).toLocaleString();
  }

  render() {
    const isPageChanging = this.accountScreenData.isLoading;
    const isAddressesLoadingInitially = this.accountScreenData.accounts.length === 0;
    const totalPages = Math.ceil(TOTAL_ADDRESSES_COUNT / this.accountScreenData.addressesPerPage);
    const isSelectedWalletOnPage = this.accountScreenData.accounts.some(
      accountDerivation => accountDerivation.index === this.selectedIndex,
    );
    const isAccessWalletDisabled = !isSelectedWalletOnPage && !isPageChanging;
    const lastIndexOfPage = this.accountScreenData.startIndex + this.accountScreenData.addressesPerPage;
    const isSingleDigitIndex = lastIndexOfPage <= 10;
    const isIndexBelowOneHundred = !isSingleDigitIndex && lastIndexOfPage < 100;
    const isIndexInTheHundreds = !isIndexBelowOneHundred && !isSingleDigitIndex && lastIndexOfPage < 1000;
    const isIndexInTheThousands = lastIndexOfPage >= 1000;

    const ledgerAddressesClasses: Record<string, string> = {
      buttonTooltip: 'drt:absolute drt:top-0 drt:h-12 drt:left-0 drt:right-0',
      preloaderItem:
        'drt:h-16! drt:border drt:border-solid drt:border-transparent drt:rounded-lg! drt:flex drt:items-center drt:w-full! drt:p-4',
      preloaderItemCheckbox: 'drt:h-4! drt:mr-2 drt:min-w-4! drt:w-4! drt:rounded-full! drt:bg-neutral-700!',
      preloaderItemAddress: 'drt:w-40! drt:h-4! drt:bg-neutral-700! drt:rounded-lg! drt:mr-auto',
      preloaderItemBalance: 'drt:w-24! drt:h-4! drt:bg-neutral-700! drt:rounded-lg! drt:ml-2',
      preloaderItemIndex: classNames('drt:mr-2 drt:h-4! drt:bg-neutral-700! drt:rounded-lg!', {
        'drt:w-9!': isSingleDigitIndex,
        'drt:w-10!': isIndexBelowOneHundred,
        'drt:w-13!': isIndexInTheHundreds,
        'drt:w-17!': isIndexInTheThousands,
      }),
    };

    if (isAddressesLoadingInitially) {
      return <drt-ledger-intro isAwaiting={true} />;
    }

    return (
      <div class="ledger-addresses">
        <div class="ledger-addresses-label">Choose the wallet you want to access</div>

        <div class="ledger-addresses-wrapper">
          <div class={{ 'ledger-addresses-preloader': true, 'visible': isPageChanging }}>
            {Array.from({ length: this.accountScreenData.addressesPerPage }, () => (
              <drt-preloader
                class={classNames('ledger-addresses-preloader-item', ledgerAddressesClasses.preloaderItem)}
              >
                <drt-preloader
                  class={classNames(
                    'ledger-addresses-preloader-item-checkbox',
                    ledgerAddressesClasses.preloaderItemCheckbox,
                  )}
                />
                <drt-preloader
                  class={classNames('ledger-addresses-preloader-item-index', ledgerAddressesClasses.preloaderItemIndex)}
                />
                <drt-preloader
                  class={classNames(
                    'ledger-addresses-preloader-item-address',
                    ledgerAddressesClasses.preloaderItemAddress,
                  )}
                />
                <drt-preloader
                  class={classNames(
                    'ledger-addresses-preloader-item-balance',
                    ledgerAddressesClasses.preloaderItemBalance,
                  )}
                />
              </drt-preloader>
            ))}
          </div>

          <div class={{ 'ledger-addresses-list': true, 'visible': !isPageChanging }}>
            {this.accountScreenData.accounts.map(accountDerivation => (
              <div
                class={{
                  'ledger-addresses-list-item': true,
                  'checked': accountDerivation.index === this.selectedIndex,
                }}
                onClick={this.handleSelectAccount(accountDerivation.index)}
              >
                <div
                  class={{
                    'ledger-addresses-list-item-checkbox': true,
                    'checked': accountDerivation.index === this.selectedIndex,
                  }}
                />
                <div
                  class={{
                    'ledger-addresses-list-item-index': true,
                    'checked': accountDerivation.index === this.selectedIndex,
                    'narrow': isSingleDigitIndex,
                    'middle': isIndexBelowOneHundred,
                    'larger': isIndexInTheHundreds,
                    'largest': isIndexInTheThousands,
                  }}
                >
                  #{this.processLedgerAddressIndex(accountDerivation)}
                </div>

                <drt-trim-text text={accountDerivation.address} class="ledger-addresses-list-item-address" />
                <div class="ledger-addresses-list-item-balance">{accountDerivation.usdValue}</div>
              </div>
            ))}
          </div>
        </div>

        <drt-pagination
          totalPages={totalPages}
          isDisabled={isPageChanging}
          onPageChange={(event: CustomEvent) => this.handlePageChange(event)}
          currentPage={this.accountScreenData.startIndex / this.accountScreenData.addressesPerPage + 1}
        />

        <div class="ledger-addresses-button-wrapper">
          {isAccessWalletDisabled && (
            <div class="ledger-addresses-button-tooltip-wrapper">
              <drt-tooltip
                trigger={
                  <div
                    class={{ 'ledger-addresses-button-tooltip': true, [ledgerAddressesClasses.buttonTooltip]: true }}
                  />
                }
              >
                You have to select a wallet from the list that you want to access.
              </drt-tooltip>
            </div>
          )}

          <button
            data-testid={DataTestIdsEnum.confirmBtn}
            onClick={this.handleAccessWallet.bind(this)}
            class={{ 'ledger-addresses-button': true, 'loading': isPageChanging, 'disabled': isAccessWalletDisabled }}
          >
            <span class="ledger-addresses-button-label">{isPageChanging ? 'Loading Wallets...' : 'Access Wallet'}</span>
            {isPageChanging && <drt-spinner-icon />}
          </button>
        </div>
      </div>
    );
  }
}
