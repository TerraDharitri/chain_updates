import { sanitizeCallbackUrl } from '../sanitizeCallbackUrl';

describe('sanitizeLoginCallbackUrl tests', () => {
  test('Keeps original URL unchanged', () => {
    const result = sanitizeCallbackUrl('https://wallet.dharitri.org');
    expect(result).toEqual('https://wallet.dharitri.org');
  });
  test('Keeps the path and hash of the URL', () => {
    const result = sanitizeCallbackUrl(
      'https://wallet.dharitri.org/path#hash'
    );
    expect(result).toEqual('https://wallet.dharitri.org/path#hash');
  });
  test('remove "address" item from query params when this is the only query item', () => {
    const result = sanitizeCallbackUrl(
      'https://localhost:3000/feed?address=drt1-address'
    );
    expect(result).toEqual('https://localhost:3000/feed');
  });
  test('remove only "address" item from query params when many query params exits on the URL', () => {
    const result = sanitizeCallbackUrl(
      'https://localhost:3000/feed?address=drt1-address&qp1=some-val&qp2=some-val-2'
    );
    expect(result).toEqual(
      'https://localhost:3000/feed?qp1=some-val&qp2=some-val-2'
    );
  });
  test('remove only all vulnerable items from query params', () => {
    const result = sanitizeCallbackUrl(
      'https://localhost:3000/feed?address=drt1-address&qp1=some-val&vulnerableItem1=vi1&vulnerableItem2=vi2&abc=123',
      ['address', 'vulnerableItem1', 'vulnerableItem2']
    );
    expect(result).toEqual('https://localhost:3000/feed?qp1=some-val&abc=123');
  });
  test('allow sanitize vscode origin', () => {
    const result = sanitizeCallbackUrl('vscode://numbat.vscode-numbat-ide');
    expect(result).toEqual('vscode://numbat.vscode-numbat-ide');
  });
});
