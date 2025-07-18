import { CacheService } from "@terradharitri/sdk-nestjs-cache";
import { TokenUtils } from "@terradharitri/sdk-nestjs-common";
import { Test } from "@nestjs/testing";
import { ApiConfigService } from "src/common/api-config/api.config.service";
import { AssetsService } from "src/common/assets/assets.service";
import { IndexerService } from "src/common/indexer/indexer.service";
import { PersistenceService } from "src/common/persistence/persistence.service";
import { PluginService } from "src/common/plugins/plugin.service";
import { CollectionService } from "src/endpoints/collections/collection.service";
import { CollectionFilter } from "src/endpoints/collections/entities/collection.filter";
import { NftCollection } from "src/endpoints/collections/entities/nft.collection";
import { NftCollectionDetailed } from "src/endpoints/collections/entities/nft.collection.detailed";
import { DcdtAddressService } from "src/endpoints/dcdt/dcdt.address.service";
import { DcdtService } from "src/endpoints/dcdt/dcdt.service";
import { NftType } from "src/endpoints/nfts/entities/nft.type";
import { CollectionRoles } from "src/endpoints/tokens/entities/collection.roles";
import { TokenAssetStatus } from "src/endpoints/tokens/entities/token.asset.status";
import { VmQueryService } from "src/endpoints/vm.query/vm.query.service";

describe('CollectionService', () => {
  let service: CollectionService;
  let indexerService: IndexerService;
  let dcdtAddressService: DcdtAddressService;
  let assetsService: AssetsService;

  const indexerCollectionMock = {
    _id: 'XDAY23TEAM-f7a346',
    name: 'XDAY2023TEAM',
    ticker: 'XDAY23TEAM',
    token: 'XDAY23TEAM-f7a346',
    issuer: 'drt178ah2z70a442g9hrt39w2ld67lav62jq72gzp3r9tu5egz4hr4csnlrlsl',
    currentOwner: 'drt178ah2z70a442g9hrt39w2ld67lav62jq72gzp3r9tu5egz4hr4csnlrlsl',
    numDecimals: 0,
    type: 'NonFungibleDCDT',
    timestamp: 1696923978,
    ownersHistory: [
      {
        address: 'drt178ah2z70a442g9hrt39w2ld67lav62jq72gzp3r9tu5egz4hr4csnlrlsl',
        timestamp: 1696923978,
      },
    ],
    properties: {
      canMint: false,
      canBurn: false,
      canUpgrade: true,
      canTransferNFTCreateRole: true,
      canAddSpecialRoles: true,
      canPause: true,
      canFreeze: true,
      canWipe: true,
      canChangeOwner: false,
      canCreateMultiShard: false,
    },
    nft_hasTraitSummary: true,
    roles: {
      DCDTRoleNFTCreate: ['drt178ah2z70a442g9hrt39w2ld67lav62jq72gzp3r9tu5egz4hr4csnlrlsl'],
      DCDTTransferRole: ['drt178ah2z70a442g9hrt39w2ld67lav62jq72gzp3r9tu5egz4hr4csnlrlsl'],
      DCDTRoleNFTBurn: ['drt178ah2z70a442g9hrt39w2ld67lav62jq72gzp3r9tu5egz4hr4csnlrlsl'],
    },
    nft_hasRarities: false,
    api_holderCount: 101,
    api_isVerified: true,
    api_nftCount: 131,
  };

  beforeEach(async () => {
    const moduleRef = await Test.createTestingModule({
      providers: [
        CollectionService,
        {
          provide: IndexerService,
          useValue: {
            getCollection: jest.fn(),
            getNftCollections: jest.fn(),
            getNftCollectionsByIds: jest.fn(),
            getNftCollectionCount: jest.fn(),
            getCollectionsForAddress: jest.fn(),
          },
        },
        {
          provide: DcdtService,
          useValue: {
            getCollectionProperties: jest.fn(),
            getDcdtTokenProperties: jest.fn(),
          },
        },
        {
          provide: DcdtAddressService,
          useValue: {
            getCollectionsForAddress: jest.fn(),
            getCollectionCountForAddressFromElastic: jest.fn(),
          },
        },
        {
          provide: VmQueryService,
          useValue: {
            vmQuery: jest.fn(),
          },
        },
        {
          provide: CacheService,
          useValue: {
            batchApplyAll: jest.fn(),
          },
        },
        {
          provide: CacheService,
          useValue: {
            get: jest.fn(),
            getOrSet: jest.fn(),
            batchGetAll: jest.fn(),
            batchApplyAll: jest.fn(),
          },
        },
        {
          provide: PersistenceService,
          useValue: {
            getCollectionTraits: jest.fn(),
          },
        },
        {
          provide: ApiConfigService,
          useValue:
          {
            getCollectionPropertiesFromGateway: jest.fn(),
            getIsIndexerV3FlagActive: jest.fn(),
            getDcdtContractAddress: jest.fn(),
          },
        },
        {
          provide: AssetsService,
          useValue:
          {
            getTokenAssets: jest.fn(),
            getCollectionRanks: jest.fn(),
          },
        },
        {
          provide: PluginService,
          useValue:
          {
            processCollections: jest.fn(),
          },
        },
      ],
    }).compile();

    service = moduleRef.get<CollectionService>(CollectionService);
    indexerService = moduleRef.get<IndexerService>(IndexerService);
    dcdtAddressService = moduleRef.get<DcdtAddressService>(DcdtAddressService);
    assetsService = moduleRef.get<AssetsService>(AssetsService);
  });

  it('service should be defined', () => {
    expect(service).toBeDefined();
  });

  describe('getNftCollectionCount', () => {
    it('should return NFT collection count', async () => {
      const filter = new CollectionFilter();
      jest.spyOn(indexerService, 'getNftCollectionCount').mockResolvedValue(100);

      const result = await service.getNftCollectionCount(filter);

      expect(result).toStrictEqual(100);
      expect(indexerService.getNftCollectionCount).toHaveBeenCalledWith(filter);
      expect(indexerService.getNftCollectionCount).toHaveBeenCalledTimes(1);
    });

    it('should return NFT collection count for a given owner', async () => {
      const filter = new CollectionFilter();
      filter.owner = 'drt1qga7ze0l03chfgru0a32wxqf2226nzrxnyhzer9lmudqhjgy7ycq0wn4su';
      jest.spyOn(indexerService, 'getNftCollectionCount').mockResolvedValue(5);

      const result = await service.getNftCollectionCount(filter);

      expect(result).toStrictEqual(5);
      expect(indexerService.getNftCollectionCount).toHaveBeenCalledWith(filter);
      expect(indexerService.getNftCollectionCount).toHaveBeenCalledTimes(1);
    });
  });

  describe('getCollectionCountForAddressWithRoles', () => {
    it('should return collection count for address with roles', async () => {
      const address = 'drt1qga7ze0l03chfgru0a32wxqf2226nzrxnyhzer9lmudqhjgy7ycq0wn4su';
      const filter = new CollectionFilter();

      jest.spyOn(dcdtAddressService, 'getCollectionCountForAddressFromElastic').mockResolvedValue(100);

      const result = await service.getCollectionCountForAddressWithRoles(address, filter);

      expect(result).toStrictEqual(100);
      expect(dcdtAddressService.getCollectionCountForAddressFromElastic).toHaveBeenCalledWith(address, filter);
      expect(dcdtAddressService.getCollectionCountForAddressFromElastic).toHaveBeenCalledTimes(1);
    });
  });

  describe('isCollection', () => {
    it('should return true if given collection identifier is a collections', async () => {
      const identifier = 'XDAY23TEAM-f7a346';
      jest.spyOn(indexerService, 'getCollection').mockResolvedValue(indexerCollectionMock);

      const result = await service.isCollection(identifier);

      expect(result).toBeTruthy();
      expect(indexerService.getCollection).toHaveBeenCalledWith(identifier);
      expect(indexerService.getCollection).toHaveBeenCalledTimes(1);
    });
  });

  describe('getCollection', () => {
    const propertiesToCollectionsMock: NftCollection = {
      collection: 'XDAY23TEAM-f7a346',
      type: NftType.NonFungibleDCDT,
      subType: undefined,
      name: 'xPortalAchievements',
      ticker: 'XDAY23TEAM',
      owner: 'drt1lpc6wjh2hav6q50p8y6a44r2lhtnseqksygakjfgep6c9uduchkqut4le4',
      timestamp: 0,
      canFreeze: true,
      canWipe: true,
      canPause: true,
      canTransferNftCreateRole: true,
      canChangeOwner: false,
      canUpgrade: false,
      canAddSpecialRoles: false,
      decimals: undefined,
      assets: {
        website: 'https://xday.com',
        description:
          'Test description.',
        status: TokenAssetStatus.active,
        pngUrl: 'https://media.dharitri.org/tokens/asset/XDAY23TEAM-f7a346/logo.png',
        name: '',
        svgUrl: 'https://media.dharitri.org/tokens/asset/XDAY23TEAM-f7a346/logo.svg',
        extraTokens: [''],
        ledgerSignature: '',
        priceSource: undefined,
        preferredRankAlgorithm: undefined,
        lockedAccounts: undefined,
      },
      scamInfo: undefined,
      traits: [],
      auctionStats: undefined,
      isVerified: undefined,
      holderCount: undefined,
      nftCount: undefined,
    };

    it('should return collection details for a given collection identifier', async () => {
      const identifier = 'XDAY23TEAM-f7a346';

      jest.spyOn(indexerService, 'getCollection').mockResolvedValue(indexerCollectionMock);
      jest.spyOn(service, 'applyPropertiesToCollections').mockResolvedValue([propertiesToCollectionsMock]);

      const result = await service.getNftCollection(identifier);

      expect(result).toBeInstanceOf(Object);
      expect(indexerService.getCollection).toHaveBeenCalledTimes(1);
      expect(indexerService.getCollection).toHaveBeenCalledWith(identifier);
      expect(service.applyPropertiesToCollections).toHaveBeenCalledWith([identifier]);
    });

    it('should return undefined if the collection is not found', async () => {
      const identifier = 'XDAY23TEAM';
      jest.spyOn(indexerService, 'getCollection').mockResolvedValue(undefined);

      const result = await service.getNftCollection(identifier);

      expect(result).toBeUndefined();
      expect(indexerService.getCollection).toHaveBeenCalledWith(identifier);
    });

    it('should return undefined if the identifier is not a valid collection', async () => {
      const identifier = 'XDAY23TEAM';
      jest.spyOn(indexerService, 'getCollection').mockResolvedValue(undefined);
      jest.spyOn(TokenUtils, 'isCollection').mockReturnValue(true);

      const result = await service.getNftCollection(identifier);

      expect(result).toBeUndefined();
    });

    it('should return undefined for an unsupported collection type', async () => {
      const identifier = 'XDAY23TEAM';
      jest.spyOn(indexerService, 'getCollection').mockResolvedValue({
        ...indexerCollectionMock,
        type: 'UnsupportedType',
      });

      const result = await service.getNftCollection(identifier);

      expect(result).toBeUndefined();
    });

    it('should return undefined if no additional properties are applied to the collection', async () => {
      const identifier = 'XDAY23TEAM';
      jest.spyOn(indexerService, 'getCollection').mockResolvedValue(indexerCollectionMock);
      jest.spyOn(service, 'applyPropertiesToCollections').mockResolvedValue([]);

      const result = await service.getNftCollection(identifier);

      expect(result).toBeUndefined();
    });

    it('should process the collection details fully', async () => {
      const identifier = 'XDAY23TEAM';
      jest.spyOn(indexerService, 'getCollection').mockResolvedValue(indexerCollectionMock);
      jest.spyOn(service, 'applyPropertiesToCollections').mockResolvedValue([propertiesToCollectionsMock]);
      const result = await service.getNftCollection(identifier);

      expect(result).toBeInstanceOf(NftCollectionDetailed);
    });
  });

  describe('getNftCollectionRolesFromGateway', () => {
    const gatewayCollectionRolesMock: CollectionRoles = {
      address: 'drt1lpc6wjh2hav6q50p8y6a44r2lhtnseqksygakjfgep6c9uduchkqut4le4',
      canCreate: false,
      canBurn: false,
      canAddQuantity: false,
      canUpdateAttributes: false,
      canAddUri: false,
      canTransfer: undefined,
      roles: ['DCDTRoleBurnForAll'],
    };

    it('should return collection roles from gateway', async () => {
      jest.spyOn(service, 'getNftCollectionRolesFromGateway').mockResolvedValue([gatewayCollectionRolesMock]);

      const results = await service.getNftCollectionRolesFromGateway(gatewayCollectionRolesMock);

      for (const result of results) {
        expect(result.address).toStrictEqual(gatewayCollectionRolesMock.address);
        expect(service.getNftCollectionRolesFromGateway).toHaveBeenCalledWith(gatewayCollectionRolesMock);
        expect(service.getNftCollectionRolesFromGateway).toHaveBeenCalledTimes(1);
      }
    });
  });

  describe('getLogoPng', () => {
    const assetsTokenMock = {
      website: 'https://xday.com',
      description:
        'Test description.',
      status: TokenAssetStatus.active,
      pngUrl: 'https://media.dharitri.org/tokens/asset/XDAY23TEAM-f7a346/logo.png',
      name: '',
      svgUrl: 'https://media.dharitri.org/tokens/asset/XDAY23TEAM-f7a346/logo.svg',
      extraTokens: [''],
      ledgerSignature: '',
      priceSource: undefined,
      preferredRankAlgorithm: undefined,
      lockedAccounts: undefined,
    };

    it('should return the PNG URL if available', async () => {
      const identifier = 'XDAY23TEAM-f7a346';
      const mockPngUrl = 'https://media.dharitri.org/tokens/asset/XDAY23TEAM-f7a346/logo.png';

      jest.spyOn(assetsService, 'getTokenAssets').mockResolvedValue(assetsTokenMock);
      const result = await service.getLogoPng(identifier);

      expect(result).toBe(mockPngUrl);
    });

    it('should return the SVG URL if available', async () => {
      const identifier = 'XDAY23TEAM-f7a346';
      const mockSvgUrl = 'https://media.dharitri.org/tokens/asset/XDAY23TEAM-f7a346/logo.svg';

      jest.spyOn(assetsService, 'getTokenAssets').mockResolvedValue(assetsTokenMock);
      const result = await service.getLogoSvg(identifier);

      expect(result).toBe(mockSvgUrl);
    });

    it('should return undefined if the PNG URL is not available', async () => {
      const identifier = 'XDAY23TEAM-f7a346';

      jest.spyOn(assetsService, 'getTokenAssets').mockResolvedValue(undefined);

      const result = await service.getLogoPng(identifier);

      expect(result).toBeUndefined();
    });

    it('should return undefined if the SVG URL is not available', async () => {
      const identifier = 'XDAY23TEAM-f7a346';

      jest.spyOn(assetsService, 'getTokenAssets').mockResolvedValue(undefined);

      const result = await service.getLogoSvg(identifier);

      expect(result).toBeUndefined();
    });
  });
});









