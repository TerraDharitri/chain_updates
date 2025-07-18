import {
  isCrossShardTransaction,
  IsCrossShardTransactionPropsType
} from '../helpers/isCrossShardTransaction';

describe('isCrossShardTransaction', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should return false if receiver and sender are in the same shard', () => {
    const props: IsCrossShardTransactionPropsType = {
      receiverAddress:
        'drt1dp96nx38nkw46urz9m4scax45rh9qc7aqctfe972vdeuq385d2zstsl0d4',
      senderShard: 1
    };

    const result = isCrossShardTransaction(props);

    expect(result).toBe(false);
  });

  it('should return true if receiver and sender are in different shards', () => {
    const props: IsCrossShardTransactionPropsType = {
      receiverAddress:
        'drt1dzmvw7gfrfv6tjxvah9rmwd45xqzs6q098925ucmsu69s48776dqzcuss5',
      senderShard: 1
    };

    const result = isCrossShardTransaction(props);

    expect(result).toBe(true);
  });

  it('should return false if senderAddress is provided and addresses are same shard', () => {
    const props: IsCrossShardTransactionPropsType = {
      receiverAddress:
        'drt17rp8l3waauynrhtw0233qf3kcwxxv9e8c0yhkk5l72jrndjm9j4qykfhpc',
      senderAddress:
        'drt1dzmvw7gfrfv6tjxvah9rmwd45xqzs6q098925ucmsu69s48776dqzcuss5'
    };

    const result = isCrossShardTransaction(props);

    expect(result).toBe(false);
  });

  it('should return true if senderAddress is provided and addresses are in different shards', () => {
    const props: IsCrossShardTransactionPropsType = {
      receiverAddress:
        'drt1dp96nx38nkw46urz9m4scax45rh9qc7aqctfe972vdeuq385d2zstsl0d4',
      senderAddress:
        'drt1dzmvw7gfrfv6tjxvah9rmwd45xqzs6q098925ucmsu69s48776dqzcuss5'
    };

    const result = isCrossShardTransaction(props);

    expect(result).toBe(true);
  });

  it('should return false if an error is thrown', () => {
    const props: IsCrossShardTransactionPropsType = {
      receiverAddress: 'invalidAddress',
      senderShard: 1
    };

    const result = isCrossShardTransaction(props);

    expect(result).toBe(false);
  });
});
