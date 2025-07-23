import { faFileAlt } from '@fortawesome/free-solid-svg-icons/faFileAlt';
import { faLock } from '@fortawesome/free-solid-svg-icons/faLock';
import { Component, h, Prop } from '@stencil/core';
import classNames from 'classnames';
import { DataTestIdsEnum } from 'constants/dataTestIds.enum';

import type { TransactionAccountType } from '../../transactions-table.type';

const transactionAccountClasses: Record<string, string> = {
  explorerLink: 'drt:text-blue-link!',
};

@Component({
  tag: 'drt-transaction-account',
  styleUrl: 'transaction-account.scss',
})
export class TransactionAccount {
  @Prop() account: TransactionAccountType;
  @Prop() class?: string;
  @Prop() dataTestId?: string;
  @Prop() scope: 'receiver' | 'sender';
  @Prop() showLockedAccounts: boolean = false;

  render() {
    const explorerLinkDataTestId =
      this.scope === 'receiver' ? DataTestIdsEnum.receiverLink : DataTestIdsEnum.senderLink;

    return (
      <div class={classNames(this.class, 'transaction-account')} data-testid={this.dataTestId}>
        {this.showLockedAccounts && this.account.isTokenLocked && (
          <drt-fa-icon icon={faLock} description={this.account.name} />
        )}

        {this.account.isContract && <drt-fa-icon icon={faFileAlt} description="Smart Contract" />}
        {this.account.showLink ? (
          <drt-explorer-link
            link={this.account.link}
            data-testid={explorerLinkDataTestId}
            class={transactionAccountClasses.explorerLink}
          >
            <span>{this.account.address}</span>
          </drt-explorer-link>
        ) : (
          <drt-transaction-account-name
            name={this.account.name}
            description={this.account.description}
            address={this.account.address}
          />
        )}
      </div>
    );
  }
}
