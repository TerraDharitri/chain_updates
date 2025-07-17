export const mockMoaSettingsService = () => ({
  getSettings: jest.fn().mockResolvedValue({}),
});

export const mockMoaEconomicsService = () => ({
  getMoaEconomics: jest.fn().mockResolvedValue({}),
});

export const mockMoaPairService = () => ({
  getMoaPair: jest.fn().mockResolvedValue({}),
  getMoaPairs: jest.fn().mockResolvedValue([]),
  getMoaPairsCount: jest.fn().mockResolvedValue(0),
});

export const mockMoaTokensService = () => ({
  getMoaTokens: jest.fn().mockResolvedValue([]),
  getMoaTokensCount: jest.fn().mockResolvedValue(0),
  getMoaTokenByIdentifier: jest.fn().mockResolvedValue({}),
});

export const mockMoaFarmsService = () => ({
  getMoaFarms: jest.fn().mockResolvedValue([]),
  getMoaFarmsCount: jest.fn().mockResolvedValue(0),
  getMoaTokenByIdentifier: jest.fn().mockResolvedValue({}),
});
