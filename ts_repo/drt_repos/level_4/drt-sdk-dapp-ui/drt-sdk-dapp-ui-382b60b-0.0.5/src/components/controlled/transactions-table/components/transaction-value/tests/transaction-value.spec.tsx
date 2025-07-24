import { h } from '@stencil/core';
import { newSpecPage } from '@stencil/core/testing';
import type { ITransactionValue } from 'components/controlled/transactions-table/transactions-table.type';

import { TransactionValue } from '../transaction-value';

describe('TransactionValue', () => {
  it('renders with minimal props', async () => {
    const value: ITransactionValue = {
      rewaLabel: '',
      link: '',
      linkText: '',
      name: '',
      ticker: '',
      valueDecimal: '',
      valueInteger: '',
    };

    const page = await newSpecPage({
      components: [TransactionValue],
      template: () => <drt-transaction-value value={value}></drt-transaction-value>,
    });

    expect(page.root).toEqualHtml(`
      <drt-transaction-value>
          <div class="transaction-value"></div>
      </drt-transaction-value>
    `);
  });

  it('renders with badge', async () => {
    const value: ITransactionValue = {
      badge: 'NFT',
      rewaLabel: '',
      link: '',
      linkText: '',
      name: '',
      ticker: '',
      valueDecimal: '',
      valueInteger: '',
    };

    const page = await newSpecPage({
      components: [TransactionValue],
      template: () => <drt-transaction-value value={value}></drt-transaction-value>,
    });

    expect(page.root).toEqualHtml(`
      <drt-transaction-value>
          <div class="transaction-value">
            <div class="badge badge-pill badge-secondary font-weight-light transaction-value-badge" data-testid="transactionNftBadge">
              NFT
            </div>
          </div>
      </drt-transaction-value>
    `);
  });

  it('renders with formatted amount', async () => {
    const value: ITransactionValue = {
      showFormattedAmount: true,
      rewaLabel: 'xREWA',
      valueDecimal: '123',
      valueInteger: '123',
      link: '',
      linkText: '',
      name: '',
      ticker: '',
    };

    const page = await newSpecPage({
      components: [TransactionValue],
      template: () => <drt-transaction-value value={value}></drt-transaction-value>,
    });

    expect(page.root).toEqualHtml(`
      <drt-transaction-value>
          <div class="transaction-value">
            <drt-format-amount class="mr-1" datatestid="transactionActionFormattedAmount" isvalid="" label="xREWA" valuedecimal="123" valueinteger="123"></drt-format-amount>
          </div>
      </drt-transaction-value>
    `);
  });

  it('renders with explorer link', async () => {
    const value: ITransactionValue = {
      rewaLabel: '',
      link: 'https://example.com',
      linkText: 'Example Link',
      name: 'Example',
      ticker: 'EX',
      valueDecimal: '',
      valueInteger: '',
    };

    const page = await newSpecPage({
      components: [TransactionValue],
      template: () => <drt-transaction-value value={value}></drt-transaction-value>,
    });

    expect(page.root).toEqualHtml(`
      <drt-transaction-value>
          <div class="transaction-value">
            <drt-explorer-link class="transaction-value-link text-truncate" link="https://example.com">
              <div class="transaction-value-content" slot="content">
                <span class="transaction-value-link-text">Example Link</span>
              </div>
            </drt-explorer-link>
          </div>
      </drt-transaction-value>
    `);
  });

  it('renders with SVG icon', async () => {
    const value: ITransactionValue = {
      rewaLabel: '',
      link: 'https://example.com',
      linkText: 'Example Link',
      svgUrl: 'https://example.com/icon.svg',
      name: 'Example Icon',
      ticker: 'EX',
      valueDecimal: '',
      valueInteger: '',
    };

    const page = await newSpecPage({
      components: [TransactionValue],
      template: () => <drt-transaction-value value={value}></drt-transaction-value>,
    });

    expect(page.root).toEqualHtml(`
      <drt-transaction-value>
          <div class="transaction-value">
            <drt-explorer-link class="transaction-value-link side-link d-flex" link="https://example.com">
              <div class="transaction-value-content" slot="content">
                <img alt="Example Icon" class="transaction-value-img" src="https://example.com/icon.svg">
                <span class="transaction-value-link-text">Example Link</span>
              </div>
            </drt-explorer-link>
          </div>
      </drt-transaction-value>
    `);
  });

  it('renders with truncated text', async () => {
    const value: ITransactionValue = {
      rewaLabel: '',
      link: 'https://example.com',
      linkText: 'Example Link',
      ticker: 'EXM',
      collection: 'EXM',
      name: 'Example',
      valueDecimal: '',
      valueInteger: '',
    };

    const page = await newSpecPage({
      components: [TransactionValue],
      template: () => <drt-transaction-value value={value}></drt-transaction-value>,
    });

    expect(page.root).toEqualHtml(`
      <drt-transaction-value>
          <div class="transaction-value">
            <drt-explorer-link class="transaction-value-link text-truncate" link="https://example.com">
              <div class="transaction-value-content" slot="content">
                <span class="transaction-value-link-text truncate">Example Link</span>
              </div>
            </drt-explorer-link>
          </div>
      </drt-transaction-value>
    `);
  });

  it('renders with titleText', async () => {
    const value: ITransactionValue = {
      rewaLabel: '',
      link: 'https://example.com',
      linkText: 'Example Link',
      name: 'Example',
      ticker: 'EX',
      valueDecimal: '',
      valueInteger: '',
      titleText: 'Title Text',
    };

    const page = await newSpecPage({
      components: [TransactionValue],
      template: () => <drt-transaction-value value={value}></drt-transaction-value>,
    });

    expect(page.root).toEqualHtml(`
      <drt-transaction-value>
          <div class="transaction-value">
            <drt-explorer-link class="transaction-value-link text-truncate" link="https://example.com">
              <div class="transaction-value-content" slot="content">
                <span class="transaction-value-link-text">Example Link</span>
              </div>
            </drt-explorer-link>
            <drt-fa-icon class="transaction-value-icon" title="Title Text"></drt-fa-icon>
          </div>
      </drt-transaction-value>
    `);
  });
});
