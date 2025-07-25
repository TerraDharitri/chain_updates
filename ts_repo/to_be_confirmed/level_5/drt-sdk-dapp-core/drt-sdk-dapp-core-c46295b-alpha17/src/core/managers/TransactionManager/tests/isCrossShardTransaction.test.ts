import { Address } from '@terradharitri/sdk-core/out';
import { getShardOfAddress } from '../helpers/getShardOfAddress';
import {
  isCrossShardTransaction,
  IsCrossShardTransactionPropsType
} from '../helpers/isCrossShardTransaction';

jest.mock('@terradharitri/sdk-core/out', () => ({
  Address: jest.fn().mockImplementation((address: string) => ({
    pubkey: jest.fn().mockReturnValue(address)
  }))
}));

jest.mock('../helpers/getShardOfAddress', () => ({
  getShardOfAddress: jest.fn()
}));

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
    (getShardOfAddress as jest.Mock).mockReturnValue(1);

    const result = isCrossShardTransaction(props);

    expect(result).toBe(false);
    expect(Address).toHaveBeenCalledWith(
      'drt1dp96nx38nkw46urz9m4scax45rh9qc7aqctfe972vdeuq385d2zstsl0d4'
    );
    expect(getShardOfAddress).toHaveBeenCalledWith(
      'drt1dp96nx38nkw46urz9m4scax45rh9qc7aqctfe972vdeuq385d2zstsl0d4'
    );
  });

  it('should return true if receiver and sender are in different shards', () => {
    const props: IsCrossShardTransactionPropsType = {
      receiverAddress:
        'drt1dzmvw7gfrfv6tjxvah9rmwd45xqzs6q098925ucmsu69s48776dqzcuss5',
      senderShard: 1
    };
    (getShardOfAddress as jest.Mock).mockReturnValue(2);

    const result = isCrossShardTransaction(props);

    expect(result).toBe(true);
    expect(Address).toHaveBeenCalledWith(
      'drt1dzmvw7gfrfv6tjxvah9rmwd45xqzs6q098925ucmsu69s48776dqzcuss5'
    );
    expect(getShardOfAddress).toHaveBeenCalledWith(
      'drt1dzmvw7gfrfv6tjxvah9rmwd45xqzs6q098925ucmsu69s48776dqzcuss5'
    );
  });

  it('should return false if senderAddress is provided and addresses are same shard', () => {
    const props: IsCrossShardTransactionPropsType = {
      receiverAddress:
        'drt17rp8l3waauynrhtw0233qf3kcwxxv9e8c0yhkk5l72jrndjm9j4qykfhpc',
      senderAddress:
        'drt1dzmvw7gfrfv6tjxvah9rmwd45xqzs6q098925ucmsu69s48776dqzcuss5'
    };
    (getShardOfAddress as jest.Mock).mockReturnValue(2);

    const result = isCrossShardTransaction(props);

    expect(result).toBe(false);
    expect(Address).toHaveBeenCalledWith(
      'drt17rp8l3waauynrhtw0233qf3kcwxxv9e8c0yhkk5l72jrndjm9j4qykfhpc'
    );
    expect(Address).toHaveBeenCalledWith(
      'drt17rp8l3waauynrhtw0233qf3kcwxxv9e8c0yhkk5l72jrndjm9j4qykfhpc'
    );
  });

  it('should return true if senderAddress is provided and addresses are in different shards', () => {
    const props: IsCrossShardTransactionPropsType = {
      receiverAddress:
        'drt1dp96nx38nkw46urz9m4scax45rh9qc7aqctfe972vdeuq385d2zstsl0d4',
      senderAddress:
        'drt1dzmvw7gfrfv6tjxvah9rmwd45xqzs6q098925ucmsu69s48776dqzcuss5'
    };
    (getShardOfAddress as jest.Mock).mockImplementation((address: string) => {
      if (
        address ===
        'drt1dp96nx38nkw46urz9m4scax45rh9qc7aqctfe972vdeuq385d2zstsl0d4'
      ) {
        return 1;
      }
      if (
        address ===
        'drt1dzmvw7gfrfv6tjxvah9rmwd45xqzs6q098925ucmsu69s48776dqzcuss5'
      ) {
        return 2;
      }
      return 0;
    });

    const result = isCrossShardTransaction(props);

    expect(result).toBe(true);
    expect(Address).toHaveBeenCalledWith(
      'drt1dp96nx38nkw46urz9m4scax45rh9qc7aqctfe972vdeuq385d2zstsl0d4'
    );
    expect(Address).toHaveBeenCalledWith(
      'drt1dp96nx38nkw46urz9m4scax45rh9qc7aqctfe972vdeuq385d2zstsl0d4'
    );
    expect(getShardOfAddress).toHaveBeenCalledWith(
      'drt1dp96nx38nkw46urz9m4scax45rh9qc7aqctfe972vdeuq385d2zstsl0d4'
    );
    expect(getShardOfAddress).toHaveBeenCalledWith(
      'drt1dzmvw7gfrfv6tjxvah9rmwd45xqzs6q098925ucmsu69s48776dqzcuss5'
    );
  });

  it('should return false if an error is thrown', () => {
    const props: IsCrossShardTransactionPropsType = {
      receiverAddress: 'invalidAddress',
      senderShard: 1
    };
    (Address as unknown as jest.Mock).mockImplementation(() => {
      throw new Error('Invalid address');
    });

    const result = isCrossShardTransaction(props);

    expect(result).toBe(false);
  });
});
