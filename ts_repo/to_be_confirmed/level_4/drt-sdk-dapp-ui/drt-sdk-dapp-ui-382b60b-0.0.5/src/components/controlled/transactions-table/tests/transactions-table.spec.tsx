import { faCheck } from '@fortawesome/free-solid-svg-icons';
import { h } from '@stencil/core';
import { newSpecPage } from '@stencil/core/testing';

import { TransactionsTable } from '../transactions-table';
import type { ITransactionsTableRow } from '../transactions-table.type';

describe('TransactionsTable', () => {
  const mockTransactions: ITransactionsTableRow[] = [
    {
      age: { timeAgo: '5 minutes ago', tooltip: 'Feb 7, 2025, 4:55 PM' },
      direction: 'in',
      method: { name: 'transfer', actionDescription: 'Token transfer' },
      iconInfo: { icon: faCheck, tooltip: 'Successful transaction' },
      link: '/transactions/hash1',
      receiver: {
        address: 'drt1qqqqqqqqqqqqqpgqp699jngundfqw07d8jzkepucvpzush6k3wvqeyzkqc',
        description: 'Receiver 1 Description',
        isContract: false,
        isTokenLocked: false,
        link: '/accounts/drt1qqqqqqqqqqqqqpgqp699jngundfqw07d8jzkepucvpzush6k3wvqeyzkqc',
        name: 'Receiver 1',
        shard: '0',
        shardLink: '/shard/0',
        showLink: true,
      },
      sender: {
        address: 'drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c',
        description: 'Sender 1 Description',
        isContract: false,
        isTokenLocked: false,
        link: '/accounts/drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c',
        name: 'Sender 1',
        shard: '1',
        shardLink: '/shard/1',
        showLink: true,
      },
      txHash: 'hash1',
      value: {
        badge: 'REWA',
        collection: 'native',
        rewaLabel: '100 REWA',
        link: '/token/REWA',
        linkText: 'REWA',
        name: 'DharitrI eGold',
        showFormattedAmount: true,
        svgUrl: '/assets/tokens/rewa.svg',
        ticker: 'REWA',
        titleText: '100 REWA',
        valueDecimal: '000000000000000000',
        valueInteger: '100',
      },
    },
    {
      age: { timeAgo: '10 minutes ago', tooltip: 'Feb 7, 2025, 4:50 PM' },
      direction: 'out',
      method: { name: 'dcdt_transfer', actionDescription: 'DCDT Token transfer' },
      iconInfo: { icon: faCheck, tooltip: 'Successful transaction' },
      link: '/transactions/hash2',
      receiver: {
        address: 'drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf',
        description: 'Receiver 2 Description',
        isContract: true,
        isTokenLocked: false,
        link: '/accounts/drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf',
        name: 'Receiver 2',
        shard: '2',
        shardLink: '/shard/2',
        showLink: true,
      },
      sender: {
        address: 'drt1k2s324ww2g0yj38qn2ch2jwctdy8mnfxep94q9arncc6xecg3xaq889n6e',
        description: 'Sender 2 Description',
        isContract: false,
        isTokenLocked: true,
        link: '/accounts/drt1k2s324ww2g0yj38qn2ch2jwctdy8mnfxep94q9arncc6xecg3xaq889n6e',
        name: 'Sender 2',
        shard: '1',
        shardLink: '/shard/1',
        showLink: true,
      },
      txHash: 'hash2',
      value: {
        badge: 'MOA',
        collection: 'tokens',
        rewaLabel: '500 MOA',
        link: '/token/MOA-455c57',
        linkText: 'MOA',
        name: 'MOA Token',
        showFormattedAmount: true,
        svgUrl: '/assets/tokens/moa.svg',
        ticker: 'MOA',
        titleText: '500 MOA',
        valueDecimal: '000000000000000000',
        valueInteger: '500',
      },
    },
  ];

  it('renders with transactions', async () => {
    const page = await newSpecPage({
      components: [TransactionsTable],
      template: () => <drt-transactions-table transactions={mockTransactions}></drt-transactions-table>,
    });

    expect(page.root).toBeTruthy();
    expect(page.root.querySelectorAll('tbody tr').length).toBe(2); // 2 mock transactions
  });

  it('applies custom class', async () => {
    const page = await newSpecPage({
      components: [TransactionsTable],
      template: () => <drt-transactions-table class="custom-class" transactions={mockTransactions}></drt-transactions-table>,
    });

    expect(page.root.querySelector('table').classList.contains('custom-class')).toBeTruthy();
  });

  it('renders transaction details correctly', async () => {
    const page = await newSpecPage({
      components: [TransactionsTable],
      template: () => <drt-transactions-table transactions={mockTransactions}></drt-transactions-table>,
    });

    const rows = page.root.querySelectorAll('tbody tr');
    expect(rows.length).toBe(2);

    rows.forEach((row, index) => {
      expect(row.querySelector('drt-transaction-hash')).toBeTruthy();
      expect(row.querySelector('drt-transaction-age')).toBeTruthy();
      expect(row.querySelector('drt-transaction-shards')).toBeTruthy();
      expect(row.querySelector('drt-transaction-account[scope="sender"]')).toBeTruthy();
      expect(row.querySelector('drt-transaction-account[scope="receiver"]')).toBeTruthy();
      expect(row.querySelector('drt-transaction-method')).toBeTruthy();
      expect(row.querySelector('drt-transaction-value')).toBeTruthy();

      // Check some specific values
      expect(row.querySelector('drt-transaction-age').getAttribute('age')).toBe(mockTransactions[index].age.timeAgo);
      expect(row.querySelector('drt-transaction-method').getAttribute('method')).toBe(mockTransactions[index].method.name);
    });
  });
});
