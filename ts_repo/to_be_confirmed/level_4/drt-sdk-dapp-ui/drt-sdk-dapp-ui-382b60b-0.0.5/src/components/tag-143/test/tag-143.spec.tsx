import { newSpecPage } from '@stencil/core/testing';
import { Tag143 } from '../tag-143';

describe('tag-143', () => {
  it('renders', async () => {
    const page = await newSpecPage({
      components: [Tag143],
      html: `<tag-143></tag-143>`,
    });
    expect(page.root).toEqualHtml(`
      <tag-143>
        <mock:shadow-root>
          <slot></slot>
        </mock:shadow-root>
      </tag-143>
    `);
  });
});
