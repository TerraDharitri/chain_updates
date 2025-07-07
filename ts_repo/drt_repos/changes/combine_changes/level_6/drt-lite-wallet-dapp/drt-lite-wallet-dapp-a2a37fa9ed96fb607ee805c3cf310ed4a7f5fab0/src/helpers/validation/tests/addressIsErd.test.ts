import { addressIsErd } from '../addressIsErd';

describe('addressIsErd tests', () => {
  it('returns true for valid DharitrI addresses starting with erd', () => {
    const validAddress =
      'drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf';
    expect(addressIsErd(validAddress)).toBe(true);
  });

  it('returns false for invalid addresses', () => {
    expect(addressIsErd('invalid-address')).toBe(false);
    expect(addressIsErd('')).toBe(false);
    expect(addressIsErd(null as any)).toBe(false);
    expect(addressIsErd(undefined as any)).toBe(false);
  });

  it('returns false for Ethereum addresses', () => {
    const ethAddress = '0x71C7656EC7ab88b098defB751B7401B5f6d8976F';
    expect(addressIsErd(ethAddress)).toBe(false);
  });

  it('returns false for addresses that are similar to bech32 but not DharitrI', () => {
    const nonDrtAddress =
      'xrd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th';
    expect(addressIsErd(nonDrtAddress)).toBe(false);
  });
});
