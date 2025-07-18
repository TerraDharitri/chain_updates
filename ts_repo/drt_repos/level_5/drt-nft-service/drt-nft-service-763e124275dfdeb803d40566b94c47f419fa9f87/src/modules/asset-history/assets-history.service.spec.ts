import { Test } from '@nestjs/testing';
import { DrtElasticService } from 'src/common';
import { DrtElasticServiceMock } from 'src/common/services/drt-communication/drt-elastic.service.mock';
import { AssetsHistoryService } from './assets-history.service';

describe.skip('AssetsHistoryService', () => {
  let service: AssetsHistoryService;
  const DrtElasticServiceProvider = {
    provide: DrtElasticService,
    useClass: DrtElasticServiceMock,
  };
  beforeAll(async () => {
    const moduleRef = await Test.createTestingModule({
      providers: [DrtElasticServiceProvider, AssetsHistoryService],
    }).compile();

    service = moduleRef.get<AssetsHistoryService>(AssetsHistoryService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });
});
