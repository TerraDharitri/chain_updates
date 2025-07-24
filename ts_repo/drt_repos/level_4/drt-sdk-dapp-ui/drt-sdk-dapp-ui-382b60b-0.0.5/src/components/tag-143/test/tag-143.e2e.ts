import { newE2EPage } from '@stencil/core/testing';

describe('tag-143', () => {
  it('renders', async () => {
    const page = await newE2EPage();
    await page.setContent('<tag-143></tag-143>');

    const element = await page.find('tag-143');
    expect(element).toHaveClass('hydrated');
  });
});
