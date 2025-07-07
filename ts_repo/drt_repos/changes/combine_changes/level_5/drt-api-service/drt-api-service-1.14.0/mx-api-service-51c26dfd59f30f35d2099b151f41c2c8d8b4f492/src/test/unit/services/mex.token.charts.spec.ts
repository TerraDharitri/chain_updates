import { Test } from "@nestjs/testing";
import { MoaTokenChartsService } from "src/endpoints/moa/moa.token.charts.service";
import { GraphQlService } from "src/common/graphql/graphql.service";
import { MoaTokenChart } from "src/endpoints/moa/entities/moa.token.chart";
import { MoaTokenService } from "src/endpoints/moa/moa.token.service";
import { MoaToken } from "src/endpoints/moa/entities/moa.token";
import { CacheService } from "@terradharitri/sdk-nestjs-cache";

describe('MoaTokenChartsService', () => {
  let moaTokenChartsService: MoaTokenChartsService;
  let graphQlService: GraphQlService;
  let moaTokenService: MoaTokenService;

  beforeEach(async () => {
    const moduleRef = await Test.createTestingModule({
      providers: [
        MoaTokenChartsService,
        {
          provide: GraphQlService,
          useValue: {
            getExchangeServiceData: jest.fn(),
          },
        },
        {
          provide: MoaTokenService,
          useValue: {
            getMoaTokenByIdentifier: jest.fn(),
          },
        },
        {
          provide: CacheService,
          useValue: {
            getOrSet: jest.fn(),
          },
        },
      ],
    }).compile();

    moaTokenChartsService = moduleRef.get<MoaTokenChartsService>(MoaTokenChartsService);
    graphQlService = moduleRef.get<GraphQlService>(GraphQlService);
    moaTokenService = moduleRef.get<MoaTokenService>(MoaTokenService);
  });

  it('service should be defined', () => {
    expect(moaTokenChartsService).toBeDefined();
  });

  describe('getTokenPricesHourResolutionRaw', () => {
    it('should return an array of MoaTokenChart when data is available', async () => {
      const mockToken: MoaToken = { id: 'TOKEN-123456', symbol: 'TEST', name: 'Test Token' } as MoaToken;
      const mockData = {
        values24h: [
          { timestamp: '2023-05-08 10:00:00', value: '1.5' },
          { timestamp: '2023-05-08 11:00:00', value: '1.6' },
        ],
      };

      jest.spyOn(graphQlService, 'getExchangeServiceData').mockResolvedValue(mockData);
      jest.spyOn(moaTokenService, 'getMoaTokenByIdentifier').mockResolvedValue(mockToken);
      jest.spyOn(moaTokenChartsService as any, 'isMoaToken').mockReturnValue(true);

      const result = await moaTokenChartsService.getTokenPricesHourResolutionRaw('TOKEN-123456');

      if (result) {
        expect(result).toHaveLength(2);
        expect(result[0]).toBeInstanceOf(MoaTokenChart);
        expect(result[0].timestamp).toBe(Math.floor(new Date('2023-05-08 10:00:00').getTime() / 1000));
        expect(result[0].value).toBe(1.5);
      }
    });

    it('should return an empty array when no data is available', async () => {
      jest.spyOn(graphQlService, 'getExchangeServiceData').mockResolvedValue({});
      jest.spyOn(moaTokenChartsService as any, 'isMoaToken').mockReturnValue(true);

      const result = await moaTokenChartsService.getTokenPricesHourResolutionRaw('TOKEN-123456');

      expect(result).toEqual([]);
    });
  });

  describe('getTokenPricesDayResolutionRaw', () => {
    it('should return an array of MoaTokenChart when data is available', async () => {
      const mockToken: MoaToken = { id: 'TOKEN-123456', symbol: 'TEST', name: 'Test Token' } as MoaToken;

      const mockData = {
        latestCompleteValues: [
          { timestamp: '2023-05-01 00:00:00', value: '1.5' },
          { timestamp: '2023-05-02 00:00:00', value: '1.6' },
        ],
      };

      jest.spyOn(graphQlService, 'getExchangeServiceData').mockResolvedValue(mockData);
      jest.spyOn(moaTokenService, 'getMoaTokenByIdentifier').mockResolvedValue(mockToken);
      jest.spyOn(moaTokenChartsService as any, 'isMoaToken').mockReturnValue(true);

      const result = await moaTokenChartsService.getTokenPricesDayResolutionRaw('TOKEN-123456');

      if (result) {
        expect(result).toHaveLength(2);
        expect(result[0]).toBeInstanceOf(MoaTokenChart);
        expect(result[0].timestamp).toBe(Math.floor(new Date('2023-05-01 00:00:00').getTime() / 1000));
        expect(result[0].value).toBe(1.5);
      }
    });

    it('should return an empty array when no data is available', async () => {
      jest.spyOn(graphQlService, 'getExchangeServiceData').mockResolvedValue({});
      jest.spyOn(moaTokenChartsService as any, 'isMoaToken').mockReturnValue(true);
      const result = await moaTokenChartsService.getTokenPricesDayResolutionRaw('TOKEN-123456');

      expect(result).toEqual([]);
    });
  });

  describe('convertToMoaTokenChart', () => {
    it('should correctly convert data to MoaTokenChart array', () => {
      const inputData = [
        { timestamp: '2023-05-08 10:00:00', value: '1.5' },
        { timestamp: '2023-05-08 11:00:00', value: '1.6' },
      ];

      const result = moaTokenChartsService['convertToMoaTokenChart'](inputData);

      expect(result).toHaveLength(2);
      expect(result[0]).toBeInstanceOf(MoaTokenChart);
      expect(result[0].timestamp).toBe(Math.floor(new Date('2023-05-08 10:00:00').getTime() / 1000));
      expect(result[0].value).toBe(1.5);
      expect(result[1].timestamp).toBe(Math.floor(new Date('2023-05-08 11:00:00').getTime() / 1000));
      expect(result[1].value).toBe(1.6);
    });
  });
});
