import { checkIsValidSender } from '../helpers/checkIsValidSender';

describe('checkIsValidSender tests', () => {
  test('should return true if senderAccount is null', () => {
    expect(
      checkIsValidSender(
        null,
        'drt1axhx4kenjlae6sknq7zjg2g4fvzavv979r2fg425p62wkl84avtqd4f00j'
      )
    ).toBe(true);
  });

  test('should return true if address matches senderAccount.address', () => {
    const mockAccount = {
      address: 'drt1l5xzvt9qwykn900jyxx7rueywy8nwm5sppeh4de4qukucacr2e6sh0rz42'
    };

    expect(
      checkIsValidSender(
        mockAccount,
        'drt1l5xzvt9qwykn900jyxx7rueywy8nwm5sppeh4de4qukucacr2e6sh0rz42'
      )
    ).toBe(true);
  });

  test('should return true if address matches senderAccount.activeGuardianAddress', () => {
    const mockAccount = {
      address: 'drt1l5xzvt9qwykn900jyxx7rueywy8nwm5sppeh4de4qukucacr2e6sh0rz42',
      activeGuardianAddress:
        'drt1pu4ru4qu8jt2qmk8vh9xdp8wj234t3mhfwp4zu2fmkmvg8j4hqfqa7ealn'
    };
    expect(
      checkIsValidSender(
        mockAccount,
        'drt1pu4ru4qu8jt2qmk8vh9xdp8wj234t3mhfwp4zu2fmkmvg8j4hqfqa7ealn'
      )
    ).toBe(true);
  });

  test('should return false if address does not match any in senderAccount', () => {
    const mockAccount = {
      address: 'drt1l5xzvt9qwykn900jyxx7rueywy8nwm5sppeh4de4qukucacr2e6sh0rz42',
      activeGuardianAddress:
        'drt1pu4ru4qu8jt2qmk8vh9xdp8wj234t3mhfwp4zu2fmkmvg8j4hqfqa7ealn'
    };
    expect(
      checkIsValidSender(
        mockAccount,
        'drt1axhx4kenjlae6sknq7zjg2g4fvzavv979r2fg425p62wkl84avtqd4f00j'
      )
    ).toBe(false);
  });
});
