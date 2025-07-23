import type { EventEmitter } from '@stencil/core';
import { Component, Event, h, Prop } from '@stencil/core';
import classNames from 'classnames';
import { DataTestIdsEnum } from 'constants/dataTestIds.enum';
import type { IAddressTableData, IndexedAccountType } from 'types/address-table.types';

const TOTAL_ADDRESSES_COUNT = 5000;

@Component({
  tag: 'drt-address-table',
  styleUrl: 'address-table.scss',
  shadow: true,
})
export class AddressTable {
  @Prop() accountScreenData: IAddressTableData;
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

  private processLedgerAddressIndex(accountDerivation: IndexedAccountType) {
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

    const addressClasses: Record<string, string> = {
      pagination: 'drt:relative drt:z-1',
      buttonTooltip: 'drt:absolute drt:top-0 drt:h-12 drt:left-0 drt:right-0',
      preloaderItem:
        'drt:h-16! drt:border drt:border-solid drt:border-transparent drt:rounded-lg! drt:flex drt:items-center drt:w-full! drt:p-4',
      preloaderItemCheckbox: 'drt:h-4! drt:mr-2 drt:min-w-4! drt:w-4! drt:rounded-full! drt:bg-preloader!',
      preloaderItemAddress: 'drt:w-40! drt:h-4! drt:bg-preloader! drt:rounded-lg! drt:mr-auto',
      preloaderItemBalance: 'drt:w-24! drt:h-4! drt:bg-preloader! drt:rounded-lg! drt:ml-2',
      preloaderItemIndex: classNames('drt:mr-2 drt:h-4! drt:bg-preloader! drt:rounded-lg!', {
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
      <div class="address-table">
        <div class="address-table-label-wrapper">
          <div class="address-table-label">Choose the wallet you want to access</div>
        </div>

        <div class="address-table-wrapper">
          <div class={{ 'address-table-preloader': true, 'visible': isPageChanging }}>
            {Array.from({ length: this.accountScreenData.addressesPerPage }, () => (
              <drt-preloader class={classNames('address-table-preloader-item', addressClasses.preloaderItem)}>
                <drt-preloader
                  class={classNames('address-table-preloader-item-checkbox', addressClasses.preloaderItemCheckbox)}
                />
                <drt-preloader
                  class={classNames('address-table-preloader-item-index', addressClasses.preloaderItemIndex)}
                />
                <drt-preloader
                  class={classNames('address-table-preloader-item-address', addressClasses.preloaderItemAddress)}
                />
                <drt-preloader
                  class={classNames('address-table-preloader-item-balance', addressClasses.preloaderItemBalance)}
                />
              </drt-preloader>
            ))}
          </div>

          <div class={{ 'address-table-list': true, 'visible': !isPageChanging }}>
            {this.accountScreenData.accounts.map(accountDerivation => (
              <div
                class={{
                  'address-table-list-item': true,
                  'checked': accountDerivation.index === this.selectedIndex,
                }}
                data-testid={`${DataTestIdsEnum.addressTableItem}-${accountDerivation.address}`}
                onClick={this.handleSelectAccount(accountDerivation.index)}
              >
                <div
                  class={{
                    'address-table-list-item-checkbox': true,
                    'checked': accountDerivation.index === this.selectedIndex,
                  }}
                />
                <div
                  class={{
                    'address-table-list-item-index': true,
                    'checked': accountDerivation.index === this.selectedIndex,
                    'narrow': isSingleDigitIndex,
                    'middle': isIndexBelowOneHundred,
                    'larger': isIndexInTheHundreds,
                    'largest': isIndexInTheThousands,
                  }}
                >
                  #{this.processLedgerAddressIndex(accountDerivation)}
                </div>

                <drt-trim text={accountDerivation.address} class="address-table-list-item-address" />
                <div class="address-table-list-item-balance">{accountDerivation.usdValue}</div>
              </div>
            ))}
          </div>
        </div>

        <div class="address-table-pagination">
          <drt-pagination
            totalPages={totalPages}
            isDisabled={isPageChanging}
            class={addressClasses.pagination}
            onPageChange={(event: CustomEvent) => this.handlePageChange(event)}
            currentPage={Math.floor(this.accountScreenData.startIndex / this.accountScreenData.addressesPerPage) + 1}
          />
        </div>

        <div class="address-table-button-wrapper">
          {isAccessWalletDisabled && (
            <div class="address-table-button-tooltip-wrapper">
              <drt-tooltip
                trigger={<div class={{ 'address-table-button-tooltip': true, [addressClasses.buttonTooltip]: true }} />}
              >
                You have to select a wallet from the list that you want to access.
              </drt-tooltip>
            </div>
          )}

          <button
            data-testid={DataTestIdsEnum.confirmBtn}
            onClick={this.handleAccessWallet.bind(this)}
            class={{ 'address-table-button': true, 'loading': isPageChanging, 'disabled': isAccessWalletDisabled }}
          >
            <span class="address-table-button-label">{isPageChanging ? 'Loading Wallets...' : 'Access Wallet'}</span>
            {isPageChanging && <drt-spinner-icon />}
          </button>
        </div>
      </div>
    );
  }
}
