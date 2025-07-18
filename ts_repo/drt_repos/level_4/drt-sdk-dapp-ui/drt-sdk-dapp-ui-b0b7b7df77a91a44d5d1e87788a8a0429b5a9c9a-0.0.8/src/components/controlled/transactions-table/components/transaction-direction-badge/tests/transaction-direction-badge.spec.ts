import { newSpecPage } from '@stencil/core/testing';

import { TransactionDirectionBadge } from '../transaction-direction-badge';

describe('TransactionDirectionBadge', () => {
  it('renders with "in" direction', async () => {
    const page = await newSpecPage({
      components: [TransactionDirectionBadge],
      html: '<drt-transaction-direction-badge direction="in"></drt-transaction-direction-badge>',
    });
    expect(page.root).toEqualHtml(`
      <drt-transaction-direction-badge direction="in">
          <div class="transaction-direction-badge">
            <span class="directionBadge in">
              IN
            </span>
          </div>
      </drt-transaction-direction-badge>
    `);
  });

  it('renders with "out" direction', async () => {
    const page = await newSpecPage({
      components: [TransactionDirectionBadge],
      html: '<drt-transaction-direction-badge direction="out"></drt-transaction-direction-badge>',
    });
    expect(page.root).toEqualHtml(`
      <drt-transaction-direction-badge direction="out">
          <div class="transaction-direction-badge">
            <span class="directionBadge out">
              OUT
            </span>
          </div>
      </drt-transaction-direction-badge>
    `);
  });

  it('renders with custom direction', async () => {
    const page = await newSpecPage({
      components: [TransactionDirectionBadge],
      html: '<drt-transaction-direction-badge direction="custom"></drt-transaction-direction-badge>',
    });
    expect(page.root).toEqualHtml(`
      <drt-transaction-direction-badge direction="custom">
          <div class="transaction-direction-badge">
            <span class="directionBadge custom">
              CUSTOM
            </span>
          </div>
      </drt-transaction-direction-badge>
    `);
  });
});
