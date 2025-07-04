import { canParseAmount } from '../canParseAmount';

describe('canParseAmount tests', () => {
  it('canParseAmount returns true when no decimals are provided', () => {
    const canParse = canParseAmount('1');
    expect(canParse).toBe(true);
  });

  it('canParseAmount returns true when there is a single decimal', () => {
    const canParse = canParseAmount('1.1');
    expect(canParse).toBe(true);
  });

  it('canParseAmount returns false when more decimals are provided than requested', () => {
    const canParse = canParseAmount('1.10000', 2);
    expect(canParse).toBe(false);
  });
});
