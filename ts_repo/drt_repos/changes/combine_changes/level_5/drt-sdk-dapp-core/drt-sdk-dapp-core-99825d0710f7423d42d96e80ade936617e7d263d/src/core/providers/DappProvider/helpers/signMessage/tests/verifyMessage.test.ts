import { verifyMessage } from '../verifyMessage';

jest.mock('../getVerifier', () => ({
  getVerifier: jest.fn().mockImplementation(() => ({
    verify: jest.fn().mockImplementation(() => true)
  }))
}));

describe('Verify Message test', () => {
  const signature =
    '{"address":"drt1wh9c0sjr2xn8hzf02lwwcr4jk2s84tat9ud2kaq6zr7xzpvl9l5q6pec6c","message":"0x54455354","signature":"0xfd7578037cdaed106e04c437821828a78ac0eb93d42118b6c8a11510520400cfa4a06dfa446cd6437c0f91264675bd554cbcb2e0b08622e9f210772890f12d01","version":1,"signer":"DrtJS"}';

  it('should verify message successfully', async () => {
    const { address, isVerified, message } = await verifyMessage(signature);

    expect(address).toStrictEqual(
      'drt1wh9c0sjr2xn8hzf02lwwcr4jk2s84tat9ud2kaq6zr7xzpvl9l5q6pec6c'
    );

    expect(isVerified).toStrictEqual(true);
    expect(message).toStrictEqual('TEST');
  });
});
