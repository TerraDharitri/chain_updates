import { capitalize } from '../capitalize';

describe('capitalize tests', () => {
  const str = 'DharitrI';

  it('capitalizes the word correctly', () => {
    expect(capitalize(str)).toStrictEqual('DharitrI');
  });
});
