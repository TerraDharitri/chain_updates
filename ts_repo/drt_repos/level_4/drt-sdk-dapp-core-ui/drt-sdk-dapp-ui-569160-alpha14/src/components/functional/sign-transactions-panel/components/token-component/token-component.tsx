import { Component, h } from '@stencil/core';

import state from '../../signTransactionsPanelStore';

@Component({
  tag: 'drt-token-component',
})
export class TokenComponent {
  render() {
    const { tokenTransaction } = state;
    const { amount, identifier, usdValue } = tokenTransaction || {};

    return (
      <drt-sign-transaction-component
        header={<drt-balance-component amount={amount} ticker={identifier} usdValue={usdValue} header="You are sending"></drt-balance-component>}
      ></drt-sign-transaction-component>
    );
  }
}
