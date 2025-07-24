import { CacheService } from "@terradharitri/sdk-nestjs-cache";
import { MetricsService } from "@terradharitri/sdk-nestjs-monitoring";
import { Test } from "@nestjs/testing";
import { ApiConfigService } from "src/common/api-config/api.config.service";
import { AssetsService } from "src/common/assets/assets.service";
import { GatewayService } from "src/common/gateway/gateway.service";
import { IndexerService } from "src/common/indexer/indexer.service";
import { ProtocolService } from "src/common/protocol/protocol.service";
import { CollectionService } from "src/endpoints/collections/collection.service";
import { CollectionFilter } from "src/endpoints/collections/entities/collection.filter";
import { DcdtAddressService } from "src/endpoints/dcdt/dcdt.address.service";
import { DcdtService } from "src/endpoints/dcdt/dcdt.service";
import { NftFilter } from "src/endpoints/nfts/entities/nft.filter";
import { NftExtendedAttributesService } from "src/endpoints/nfts/nft.extendedattributes.service";

describe('DcdtAddressService', () => {
  let service: DcdtAddressService;
  let indexerService: IndexerService;
  let cacheService: CacheService;
  let metricsService: MetricsService;
  let protocolService: ProtocolService;

  beforeEach(async () => {
    const moduleRef = await Test.createTestingModule({
      providers: [
        DcdtAddressService,
        {
          provide: IndexerService,
          useValue: {
            getNftCountForAddress: jest.fn(),
            getCollectionCountForAddress: jest.fn(),
            getNftsForAddress: jest.fn(),
            getNftCollections: jest.fn(),
          },
        },
        {
          provide: ApiConfigService,
          useValue:
          {
            getExternalMediaUrl: jest.fn(),
          },
        },
        {
          provide: DcdtService,
          useValue: {
            getDcdtTokenProperties: jest.fn(),
          },
        },
        {
          provide: GatewayService,
          useValue: {
            get: jest.fn(),
            getAddressDcdtRoles: jest.fn(),
          },
        },
        {
          provide: CacheService,
          useValue: {
            get: jest.fn(),
            set: jest.fn(),
            getLocal: jest.fn(),
            setLocal: jest.fn(),
          },
        },
        {
          provide: MetricsService,
          useValue:
          {
            incrementPendingApiHit: jest.fn(),
            incrementCachedApiHit: jest.fn(),
          },
        },
        {
          provide: ProtocolService,
          useValue: {
            getSecondsRemainingUntilNextRound: jest.fn(),
          },
        },
        {
          provide: NftExtendedAttributesService,
          useValue: {
            getTags: jest.fn(),
          },
        },
        {
          provide: CollectionService,
          useValue: {
            applyPropertiesToCollections: jest.fn(),
          },
        },
        {
          provide: AssetsService,
          useValue: {
            getTokenAssets: jest.fn(),
          },
        },
      ],
    }).compile();

    service = moduleRef.get<DcdtAddressService>(DcdtAddressService);
    indexerService = moduleRef.get<IndexerService>(IndexerService);
    cacheService = moduleRef.get<CacheService>(CacheService);
    metricsService = moduleRef.get<MetricsService>(MetricsService);
    protocolService = moduleRef.get<ProtocolService>(ProtocolService);
  });

  it('service should be defined', () => {
    expect(service).toBeDefined();
  });

  describe('getNftCountForAddressFromElastic', () => {
    it('should return NFT count for a given address from Elastic', async () => {
      const address = 'drt1qga7ze0l03chfgru0a32wxqf2226nzrxnyhzer9lmudqhjgy7ycq0wn4su';
      const filter = new NftFilter();
      jest.spyOn(indexerService, 'getNftCountForAddress').mockResolvedValue(10);

      const result = await service.getNftCountForAddressFromElastic(address, filter);

      expect(result).toStrictEqual(10);
      expect(indexerService.getNftCountForAddress).toHaveBeenCalledTimes(1);
      expect(indexerService.getNftCountForAddress).toHaveBeenCalledWith(address, filter);
    });

    it('should return all NSFW NFTs count for a given address from Elastic', async () => {
      const address = 'drt1qga7ze0l03chfgru0a32wxqf2226nzrxnyhzer9lmudqhjgy7ycq0wn4su';
      const filter = new NftFilter();
      filter.isNsfw = true;
      jest.spyOn(indexerService, 'getNftCountForAddress').mockResolvedValue(1);

      const result = await service.getNftCountForAddressFromElastic(address, filter);

      expect(result).toStrictEqual(1);
      expect(indexerService.getNftCountForAddress).toHaveBeenCalledTimes(1);
      expect(indexerService.getNftCountForAddress).toHaveBeenCalledWith(address, filter);
    });

    it('should return all NFT count from a specific collection for a given address from Elastic', async () => {
      const address = 'drt1qga7ze0l03chfgru0a32wxqf2226nzrxnyhzer9lmudqhjgy7ycq0wn4su';
      const filter = new NftFilter();
      filter.collection = 'XDAY23TEAM-f7a346';
      jest.spyOn(indexerService, 'getNftCountForAddress').mockResolvedValue(20);

      const result = await service.getNftCountForAddressFromElastic(address, filter);

      expect(result).toStrictEqual(20);
      expect(indexerService.getNftCountForAddress).toHaveBeenCalledTimes(1);
      expect(indexerService.getNftCountForAddress).toHaveBeenCalledWith(address, filter);
    });

    it('should return all NFT count from specific collections for a given address from Elastic', async () => {
      const address = 'drt1qga7ze0l03chfgru0a32wxqf2226nzrxnyhzer9lmudqhjgy7ycq0wn4su';
      const filter = new NftFilter();
      filter.collections = ['XDAY23TEAM-f7a346', 'XDAY23TEAM-f11111'];
      jest.spyOn(indexerService, 'getNftCountForAddress').mockResolvedValue(25);

      const result = await service.getNftCountForAddressFromElastic(address, filter);

      expect(result).toStrictEqual(25);
      expect(indexerService.getNftCountForAddress).toHaveBeenCalledTimes(1);
      expect(indexerService.getNftCountForAddress).toHaveBeenCalledWith(address, filter);
    });
  });

  describe('getCollectionCountForAddressFromElastic', () => {
    it('should return collections count for a given address', async () => {
      const address = 'drt1qga7ze0l03chfgru0a32wxqf2226nzrxnyhzer9lmudqhjgy7ycq0wn4su';
      const filter = new CollectionFilter();
      jest.spyOn(indexerService, 'getCollectionCountForAddress').mockResolvedValue(5);

      const result = await service.getCollectionCountForAddressFromElastic(address, filter);

      expect(result).toStrictEqual(5);
      expect(indexerService.getCollectionCountForAddress).toHaveBeenCalledTimes(1);
      expect(indexerService.getCollectionCountForAddress).toHaveBeenCalledWith(address, filter);
    });

    it('should return collection count when collection filter is applied for a given address', async () => {
      const address = 'drt1qga7ze0l03chfgru0a32wxqf2226nzrxnyhzer9lmudqhjgy7ycq0wn4su';
      const filter = new CollectionFilter();
      filter.collection = 'XDAY23TEAM-f7a346';
      jest.spyOn(indexerService, 'getCollectionCountForAddress').mockResolvedValue(1);

      const result = await service.getCollectionCountForAddressFromElastic(address, filter);

      expect(result).toStrictEqual(1);
      expect(indexerService.getCollectionCountForAddress).toHaveBeenCalledTimes(1);
      expect(indexerService.getCollectionCountForAddress).toHaveBeenCalledWith(address, filter);
    });

    it('should return collection count when identifiers filter is applied for a given address', async () => {
      const address = 'drt1qga7ze0l03chfgru0a32wxqf2226nzrxnyhzer9lmudqhjgy7ycq0wn4su';
      const filter = new CollectionFilter();
      filter.identifiers = ['XDAY23TEAM-f7a346-01', 'XDAY23TEAM-f7a346-02'];
      jest.spyOn(indexerService, 'getCollectionCountForAddress').mockResolvedValue(1);

      const result = await service.getCollectionCountForAddressFromElastic(address, filter);

      expect(result).toStrictEqual(1);
      expect(indexerService.getCollectionCountForAddress).toHaveBeenCalledTimes(1);
      expect(indexerService.getCollectionCountForAddress).toHaveBeenCalledWith(address, filter);
    });
  });

  describe('DcdtAddressService - getCollectionsForAddress', () => {
    it('should return cached DCDTs for a given address', async () => {
      const address = 'drt1qga7ze0l03chfgru0a32wxqf2226nzrxnyhzer9lmudqhjgy7ycq0wn4su';
      const cachedDcdts = {
        "BUSD-19079b": {
          balance: "50850380000",
          tokenIdentifier: "BUSD-19079b",
        },
      };

      jest.spyOn(cacheService, 'getLocal').mockResolvedValue(cachedDcdts);
      jest.spyOn(metricsService, 'incrementCachedApiHit');

      const result = await service.getAllDcdtsForAddressFromGateway(address);

      expect(result).toEqual(cachedDcdts);
      expect(metricsService.incrementCachedApiHit).toHaveBeenCalledWith('Gateway.AccountDcdts');
    });

    it('should update the cache with new DCDTs data for a given address', async () => {
      const address = 'some-address';
      const ttl = 1000;

      jest.spyOn(cacheService, 'getLocal').mockResolvedValueOnce(null);
      jest.spyOn(cacheService, 'setLocal');
      jest.spyOn(protocolService, 'getSecondsRemainingUntilNextRound').mockResolvedValue(ttl);

      await service.getAllDcdtsForAddressFromGateway(address);
      const result = await service.getAllDcdtsForAddressFromGateway(address);

      expect(cacheService.setLocal).toHaveBeenCalledWith(`address:${address}:dcdts`, expect.anything(), ttl);
      expect(result).toBeDefined();
    });
  });
});
