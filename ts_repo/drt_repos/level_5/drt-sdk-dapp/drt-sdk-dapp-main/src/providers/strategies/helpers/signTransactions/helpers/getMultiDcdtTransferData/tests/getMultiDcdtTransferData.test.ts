import { getMultiDcdtTransferData } from '../getMultiDcdtTransferData';

describe('getMultiDcdtTransferData', () => {
  it('should extract transaction information', async () => {
    const data = getMultiDcdtTransferData([]);
    // Assert the result is correct based on your mock data
    expect(data).toBeTruthy();
  });
});
