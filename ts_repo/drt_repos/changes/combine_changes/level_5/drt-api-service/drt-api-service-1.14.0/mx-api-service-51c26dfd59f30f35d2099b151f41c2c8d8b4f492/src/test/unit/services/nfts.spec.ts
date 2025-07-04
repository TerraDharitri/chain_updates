import { CacheService } from "@terradharitri/sdk-nestjs-cache";
import { Test } from "@nestjs/testing";
import { ApiConfigService } from "src/common/api-config/api.config.service";
import { AssetsService } from "src/common/assets/assets.service";
import { IndexerService } from "src/common/indexer/indexer.service";
import { LockedAssetService } from "src/common/locked-asset/locked-asset.service";
import { PersistenceService } from "src/common/persistence/persistence.service";
import { PluginService } from "src/common/plugins/plugin.service";
import { DcdtAddressService } from "src/endpoints/dcdt/dcdt.address.service";
import { DcdtService } from "src/endpoints/dcdt/dcdt.service";
import { MoaTokenService } from "src/endpoints/moa/moa.token.service";
import { NftService } from "src/endpoints/nfts/nft.service";
import { NftMediaService } from "src/queue.worker/nft.worker/queue/job-services/media/nft.media.service";
import { NftMetadataService } from "src/queue.worker/nft.worker/queue/job-services/metadata/nft.metadata.service";

describe('NftService', () => {
  let service: NftService;
  let indexerService: IndexerService;

  beforeEach(async () => {
    const moduleRef = await Test.createTestingModule({
      providers: [
        NftService,
        {
          provide: ApiConfigService,
          useValue: {
            getExternalMediaUrl: jest.fn(),
            isNftExtendedAttributesEnabled: jest.fn(),
            getIsIndexerV3FlagActive: jest.fn(),
            getNftExtendedAttributesNsfwThreshold: jest.fn(),
          },
        },
        {
          provide: IndexerService,
          useValue: {
            getNfts: jest.fn(),
            getNftOwnersCount: jest.fn(),
            getNftCount: jest.fn(),
            getAccountDcdtByIdentifiers: jest.fn(),
            getAccountsDcdtByCollection: jest.fn(),
          },
        },
        {
          provide: DcdtService,
          useValue: {
            getTokenSupply: jest.fn(),
            getDcdtTokenProperties: jest.fn(),
          },
        },
        {
          provide: AssetsService,
          useValue: {
            getTokenAssets: jest.fn(),
          },
        },
        {
          provide: CacheService,
          useValue: {
            batchApplyAll: jest.fn(),
            batchApply: jest.fn(),
            getOrSet: jest.fn(),
          },
        },
        {
          provide: PluginService,
          useValue: {
            processNfts: jest.fn(),
            batchApply: jest.fn(),
          },
        },
        {
          provide: NftMetadataService,
          useValue: {
            getMetadata: jest.fn(),
          },
        },
        {
          provide: NftMediaService,
          useValue: {
            getMedia: jest.fn(),
          },
        },
        {
          provide: PersistenceService,
          useValue: {
            batchGetMedia: jest.fn(),
            batchGetMetadata: jest.fn(),
          },
        },
        {
          provide: DcdtAddressService,
          useValue: {
            getNftsForAddress: jest.fn(),
            getNftCountForAddressFromElastic: jest.fn(),
          },
        },
        {
          provide: MoaTokenService,
          useValue: {
            getMoaPrices: jest.fn(),
          },
        },
        {
          provide: LockedAssetService,
          useValue: {
            getLkmoaUnlockSchedule: jest.fn(),
            getXmoaUnlockEpoch: jest.fn(),
          },
        },
      ],
    }).compile();

    service = moduleRef.get<NftService>(NftService);
    indexerService = moduleRef.get<IndexerService>(IndexerService);
  });

  it('service should be defined', () => {
    expect(service).toBeDefined();
  });

  describe('getNftOwnersCountRaw', () => {
    it('should return NFT owner count raw', async () => {
      jest.spyOn(service, 'isNft').mockResolvedValue(true);
      jest.spyOn(indexerService, 'getNftOwnersCount').mockResolvedValue(10);

      const identifier = 'XDAY23TEAM-f7a346-01';
      const result = await service.getNftOwnersCountRaw(identifier);

      expect(result).toStrictEqual(10);
      expect(service.isNft).toHaveBeenCalledWith(identifier);
      expect(indexerService.getNftOwnersCount).toHaveBeenCalledWith(identifier);
      expect(indexerService.getNftOwnersCount).toHaveBeenCalledTimes(1);
    });

    it('should return null because test simulates that given identifier is not a valid NFT', async () => {
      jest.spyOn(service, 'isNft').mockResolvedValue(false);

      const identifier = 'XDAY23TEAM-f7a346-0102';
      const result = await service.getNftOwnersCountRaw(identifier);

      expect(result).toBeNull();
    });
  });
});
