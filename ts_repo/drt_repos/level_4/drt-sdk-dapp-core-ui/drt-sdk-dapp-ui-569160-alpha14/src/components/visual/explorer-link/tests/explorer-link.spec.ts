import { faCheck } from '@fortawesome/free-solid-svg-icons';
import { newSpecPage } from '@stencil/core/testing';

import { ExplorerLink } from '../explorer-link';

describe('ExplorerLink', () => {
  it('renders with default props', async () => {
    const page = await newSpecPage({
      components: [ExplorerLink],
      html: '<drt-explorer-link link="https://example.com"></drt-explorer-link>',
    });
    expect(page.root).toEqualHtml(`
      <drt-explorer-link link="https://example.com">
        <a href="https://example.com" target="_blank" class="drt:explorer-link" rel="noreferrer">
          <drt-fa-icon class="explorer-link-icon"></drt-fa-icon>
        </a>
      </drt-explorer-link>
    `);
  });

  it('renders with custom text', async () => {
    const page = await newSpecPage({
      components: [ExplorerLink],
      html: '<drt-explorer-link link="https://example.com" text="View on Explorer"></drt-explorer-link>',
    });
    expect(page.root).toEqualHtml(`
      <drt-explorer-link link="https://example.com" text="View on Explorer">
        <a href="https://example.com" target="_blank" class="drt:explorer-link" rel="noreferrer">
          View on Explorer
        </a>
      </drt-explorer-link>
    `);
  });

  it('renders with custom class', async () => {
    const page = await newSpecPage({
      components: [ExplorerLink],
      html: '<drt-explorer-link link="https://example.com" class="custom-class"></drt-explorer-link>',
    });
    expect(page.root).toEqualHtml(`
      <drt-explorer-link link="https://example.com" class="custom-class">
        <a href="https://example.com" target="_blank" class="drt:explorer-link custom-class" rel="noreferrer">
          <drt-fa-icon class="explorer-link-icon"></drt-fa-icon>
        </a>
      </drt-explorer-link>
    `);
  });

  it('renders with custom icon', async () => {
    const page = await newSpecPage({
      components: [ExplorerLink],
      html: '<drt-explorer-link link="https://example.com"></drt-explorer-link>',
    });

    page.root.icon = faCheck;
    await page.waitForChanges();

    expect(page.root).toEqualHtml(`
      <drt-explorer-link link="https://example.com">
        <a href="https://example.com" target="_blank" class="drt:explorer-link" rel="noreferrer">
          <drt-fa-icon class="explorer-link-icon"></drt-fa-icon>
        </a>
      </drt-explorer-link>
    `);

    jest.restoreAllMocks();
  });
});
