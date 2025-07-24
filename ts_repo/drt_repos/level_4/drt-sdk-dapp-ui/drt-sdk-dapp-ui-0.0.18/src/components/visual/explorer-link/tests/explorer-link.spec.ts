import { newSpecPage } from '@stencil/core/testing';

import { ExplorerLink } from '../explorer-link';

describe('ExplorerLink', () => {
  it('renders with default props', async () => {
    const page = await newSpecPage({
      components: [ExplorerLink],
      html: '<drt-explorer-link link="https://example.com"></drt-explorer-link>',
    });

    expect(page.root.outerHTML).toEqualHtml(`
      <drt-explorer-link link="https://example.com">
        <template shadowrootmode="open">
          <a href="https://example.com" target="_blank" class="explorer-link" rel="noreferrer">
            <drt-arrow-up-right-from-square-icon class="explorer-link-icon drt:duration-200 drt:ease-in-out drt:flex drt:hover:opacity-80 drt:justify-center drt:transition-opacity"></drt-arrow-up-right-from-square-icon>
          </a>
        </template>
      </drt-explorer-link>
    `);
  });

  it('renders with custom text', async () => {
    const page = await newSpecPage({
      components: [ExplorerLink],
      html: '<drt-explorer-link link="https://example.com">View on Explorer</drt-explorer-link>',
    });

    expect(page.root.outerHTML).toEqualHtml(`
      <drt-explorer-link link="https://example.com">
        <template shadowrootmode="open">
          <a href="https://example.com" target="_blank" class="explorer-link" rel="noreferrer">
            <slot></slot>
          </a>
        </template>
        View on Explorer
      </drt-explorer-link>
    `);
  });

  it('renders with custom class', async () => {
    const page = await newSpecPage({
      components: [ExplorerLink],
      html: '<drt-explorer-link link="https://example.com" class="custom-class"></drt-explorer-link>',
    });

    expect(page.root.outerHTML).toEqualHtml(`
      <drt-explorer-link link="https://example.com" class="custom-class">
        <template shadowrootmode="open">
          <a href="https://example.com" target="_blank" class="explorer-link custom-class" rel="noreferrer">
            <drt-arrow-up-right-from-square-icon class="explorer-link-icon drt:duration-200 drt:ease-in-out drt:flex drt:hover:opacity-80 drt:justify-center drt:transition-opacity"></drt-arrow-up-right-from-square-icon>
          </a>
        </template>
      </drt-explorer-link>
    `);
  });

  it('renders with custom icon', async () => {
    const page = await newSpecPage({
      components: [ExplorerLink],
      html: '<drt-explorer-link link="https://example.com"><drt-arrow-right-icon /></drt-explorer-link>',
    });

    expect(page.root.outerHTML).toEqualHtml(`
      <drt-explorer-link link="https://example.com">
        <template shadowrootmode="open">
          <a href="https://example.com" target="_blank" class="explorer-link" rel="noreferrer">
            <slot></slot>
          </a>
        </template>
        <drt-arrow-right-icon></drt-arrow-right-icon>
      </drt-explorer-link>
    `);

    jest.restoreAllMocks();
  });
});
