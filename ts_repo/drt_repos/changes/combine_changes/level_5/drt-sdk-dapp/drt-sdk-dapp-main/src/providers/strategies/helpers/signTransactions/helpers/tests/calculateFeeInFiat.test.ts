import { calculateFeeInFiat } from '../calculateFeeInFiat';

describe('calculateFeeInFiat tests', () => {
  it('computes correct fee in fiat', () => {
    const fee = calculateFeeInFiat({
      feeLimit: '50000000000000',
      rewaPriceInUsd: 135.78
    });
    expect(fee).toBe('≈ $0.0068');
  });
});
