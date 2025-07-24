import { forwardRef, Inject, Injectable } from "@nestjs/common";
import { ApiConfigService } from "src/common/api-config/api.config.service";
import { QueryPagination } from "src/common/entities/query.pagination";
import { TokenHelpers } from "src/utils/token.helpers";
import { Nft } from "./entities/nft";
import { NftAccount } from "./entities/nft.account";
import { NftFilter } from "./entities/nft.filter";
import { NftOwner } from "./entities/nft.owner";
import { NftType } from "./entities/nft.type";
import { NftQueryOptions } from "./entities/nft.query.options";
import { DcdtService } from "../dcdt/dcdt.service";
import { AssetsService } from "../../common/assets/assets.service";
import { NftMetadataService } from "src/queue.worker/nft.worker/queue/job-services/metadata/nft.metadata.service";
import { NftMediaService } from "src/queue.worker/nft.worker/queue/job-services/media/nft.media.service";
import { NftMedia } from "./entities/nft.media";
import { CacheInfo } from "src/utils/cache.info";
import { DcdtDataSource } from "../dcdt/entities/dcdt.data.source";
import { DcdtAddressService } from "../dcdt/dcdt.address.service";
import { PersistenceService } from "src/common/persistence/persistence.service";
import { MoaTokenService } from "../moa/moa.token.service";
import { BinaryUtils, NumberUtils, RecordUtils, BatchUtils, TokenUtils, OriginLogger } from "@terradharitri/sdk-nestjs-common";
import { ApiUtils } from "@terradharitri/sdk-nestjs-http";
import { CacheService } from "@terradharitri/sdk-nestjs-cache";
import { IndexerService } from "src/common/indexer/indexer.service";
import { LockedAssetService } from "../../common/locked-asset/locked-asset.service";
import { CollectionAccount } from "../collections/entities/collection.account";
import { NftRankAlgorithm } from "src/common/assets/entities/nft.rank.algorithm";
import { NftRarity } from "./entities/nft.rarity";
import { NftRarities } from "./entities/nft.rarities";
import { SortCollectionNfts } from "../collections/entities/sort.collection.nfts";
import { TokenAssets } from "src/common/assets/entities/token.assets";
import { ScamInfo } from "src/common/entities/scam-info.dto";
import { NftSubType } from "./entities/nft.sub.type";

@Injectable()
export class NftService {
  private readonly logger = new OriginLogger(NftService.name);
  private readonly NFT_THUMBNAIL_PREFIX: string;
  readonly DEFAULT_MEDIA: NftMedia[];

  constructor(
    private readonly apiConfigService: ApiConfigService,
    private readonly indexerService: IndexerService,
    private readonly dcdtService: DcdtService,
    private readonly assetsService: AssetsService,
    private readonly cachingService: CacheService,
    private readonly nftMetadataService: NftMetadataService,
    private readonly nftMediaService: NftMediaService,
    private readonly persistenceService: PersistenceService,
    @Inject(forwardRef(() => DcdtAddressService))
    private readonly dcdtAddressService: DcdtAddressService,
    private readonly moaTokenService: MoaTokenService,
    private readonly lockedAssetService: LockedAssetService,
  ) {
    this.NFT_THUMBNAIL_PREFIX = this.apiConfigService.getExternalMediaUrl() + '/nfts/asset';
    this.DEFAULT_MEDIA = [
      {
        url: this.nftMediaService.NFT_THUMBNAIL_DEFAULT,
        originalUrl: this.nftMediaService.NFT_THUMBNAIL_DEFAULT,
        thumbnailUrl: this.nftMediaService.NFT_THUMBNAIL_DEFAULT,
        fileType: 'image/png',
        fileSize: 29512,
      },
    ];
  }

  async getNfts(queryPagination: QueryPagination, filter: NftFilter, queryOptions?: NftQueryOptions): Promise<Nft[]> {
    const { from, size } = queryPagination;

    const nfts = await this.getNftsInternal({ from, size }, filter);

    for (const nft of nfts) {
      await this.applyAssetsAndTicker(nft);
    }

    if (queryOptions && queryOptions.withOwner) {
      const nftsIdentifiers = nfts.filter(x => x.type === NftType.NonFungibleDCDT).map(x => x.identifier);

      const accountsDcdts = await this.getAccountDcdtByIdentifiers(nftsIdentifiers, { from: 0, size: nftsIdentifiers.length });

      for (const nft of nfts) {
        if (nft.type === NftType.NonFungibleDCDT) {
          const accountDcdt = accountsDcdts.find((accountDcdt: any) => accountDcdt.identifier == nft.identifier);
          if (accountDcdt) {
            nft.owner = accountDcdt.address;
          }
        }
      }
    }

    if (queryOptions && queryOptions.withSupply) {
      const supplyNfts = nfts.filter(nft => nft.type.in(NftType.SemiFungibleDCDT, NftType.MetaDCDT));
      await this.batchApplySupply(supplyNfts);
    }

    await this.batchProcessNfts(nfts);

    for (const nft of nfts) {
      await this.applyUnlockFields(nft);
    }

    return nfts;
  }

  private async batchProcessNfts(nfts: Nft[], fields?: string[]) {
    await Promise.all([
      this.batchApplyMedia(nfts, fields),
      this.batchApplyMetadata(nfts, fields),
    ]);
  }

  private async applyNftOwner(nft: Nft): Promise<void> {
    if (nft.type === NftType.NonFungibleDCDT) {
      const accountsDcdt = await this.getAccountDcdtByIdentifier(nft.identifier);
      if (accountsDcdt.length > 0) {
        nft.owner = accountsDcdt[0].address;
      }
    }
  }

  private async batchApplySupply(nfts: Nft[], fields?: string[]) {
    if (fields && !fields.includes('supply')) {
      return;
    }

    await this.cachingService.batchApplyAll(
      nfts,
      nft => CacheInfo.TokenSupply(nft.identifier).key,
      nft => this.dcdtService.getTokenSupply(nft.identifier),
      (nft, value) => nft.supply = value.totalSupply,
      CacheInfo.TokenSupply('').ttl,
    );
  }

  private async batchApplyMedia(nfts: Nft[], fields?: string[]) {
    if (fields && !fields.includes('media')) {
      return;
    }

    await this.cachingService.batchApply(
      nfts,
      nft => CacheInfo.NftMedia(nft.identifier).key,
      async nfts => {
        const getMediaResults = await this.persistenceService.batchGetMedia(nfts.map((nft) => nft.identifier));

        return RecordUtils.mapKeys(getMediaResults, identifier => CacheInfo.NftMedia(identifier).key);
      },
      (nft, media) => nft.media = media,
      CacheInfo.NftMedia('').ttl,
    );

    for (const nft of nfts) {
      if (TokenHelpers.needsDefaultMedia(nft)) {
        nft.media = this.DEFAULT_MEDIA;
      }

      this.applyRedirectMedia(nft);
    }
  }

  private async batchApplyMetadata(nfts: Nft[], fields?: string[]) {
    if (fields && !fields.includes('metadata')) {
      return;
    }

    await this.cachingService.batchApply(
      nfts,
      nft => CacheInfo.NftMetadata(nft.identifier).key,
      async nfts => {
        const getMetadataResults = await this.persistenceService.batchGetMetadata(nfts.map((nft) => nft.identifier));

        return RecordUtils.mapKeys(getMetadataResults, identifier => CacheInfo.NftMetadata(identifier).key);
      },
      (nft, metadata) => nft.metadata = metadata,
      CacheInfo.NftMetadata('').ttl
    );
  }

  private async processNft(nft: Nft) {
    await Promise.all([
      this.applyMedia(nft),
      this.applyMetadata(nft),
    ]);

    if (TokenHelpers.needsDefaultMedia(nft)) {
      nft.media = this.DEFAULT_MEDIA;
    }
  }

  async applyAssetsAndTicker(token: Nft, fields?: string[]) {
    if (fields && fields.includesNone(['ticker', 'assets'])) {
      return;
    }

    token.assets = await this.assetsService.getTokenAssets(token.identifier) ??
      await this.assetsService.getTokenAssets(token.collection);

    if (token.assets) {
      token.ticker = token.collection.split('-')[0];
    } else {
      token.ticker = token.collection;
    }
  }

  async getSingleNft(identifier: string): Promise<Nft | undefined> {
    const nfts = await this.getNftsInternal(new QueryPagination({ from: 0, size: 1 }), new NftFilter(), identifier);

    if (!TokenUtils.isNft(identifier)) {
      return undefined;
    }

    if (nfts.length === 0) {
      return undefined;
    }

    const nft: Nft = ApiUtils.mergeObjects(new Nft(), nfts[0]);

    if (nft.identifier.toLowerCase() !== identifier.toLowerCase()) {
      return undefined;
    }

    if (nft.type && nft.type.in(
      NftType.SemiFungibleDCDT, NftType.MetaDCDT,
      NftSubType.DynamicSemiFungibleDCDT, NftSubType.DynamicMetaDCDT
    )) {
      await this.applySupply(nft);
    }

    await this.applyNftOwner(nft);

    await this.applyNftAttributes(nft);

    await this.applyAssetsAndTicker(nft);

    await this.applyUnlockFields(nft);

    await this.processNft(nft);

    return nft;
  }

  private async applyUnlockFields(nft: Nft, fields?: string[]): Promise<void> {
    if (fields && fields.includesNone(['unlockSchedule', 'unlockEpoch'])) {
      return;
    }

    if (!nft.attributes) {
      return;
    }

    try {
      nft.unlockSchedule = await this.lockedAssetService.getLkmoaUnlockSchedule(nft.identifier, nft.attributes);
    } catch (error) {
      this.logger.error(`An error occurred while applying unlock schedule for NFT with identifier '${nft.identifier}' and attributes '${nft.attributes}'`);
      this.logger.error(error);
    }

    try {
      nft.unlockEpoch = await this.lockedAssetService.getXmoaUnlockEpoch(nft.identifier, nft.attributes);
    } catch (error) {
      this.logger.error(`An error occurred while applying unlock epoch for NFT with identifier '${nft.identifier}' and attributes '${nft.attributes}'`);
      this.logger.error(error);
    }
  }

  private async applyNftAttributes(nft: Nft): Promise<void> {
    if (!nft.owner) {
      return;
    }

    const nftsForAddress = await this.dcdtAddressService.getNftsForAddress(nft.owner, new NftFilter({ identifiers: [nft.identifier] }), new QueryPagination({ from: 0, size: 1 }));
    if (nftsForAddress.length === 0) {
      return;
    }

    nft.attributes = nftsForAddress[0].attributes;
  }

  private async applyMedia(nft: Nft) {
    nft.media = await this.nftMediaService.getMedia(nft.identifier) ?? undefined;

    this.applyRedirectMedia(nft);
  }

  private async applyMetadata(nft: Nft) {
    nft.metadata = await this.nftMetadataService.getMetadata(nft) ?? undefined;
  }

  async isNft(identifier: string): Promise<boolean> {
    if (identifier.split('-').length !== 3) {
      return false;
    }

    const nfts = await this.getNftsInternal(new QueryPagination({ from: 0, size: 1 }), new NftFilter(), identifier);

    return nfts.length > 0;
  }

  async getNftOwners(identifier: string, pagination: QueryPagination): Promise<NftOwner[] | undefined> {
    const isNft = await this.isNft(identifier);
    if (!isNft) {
      return undefined;
    }

    const accountsDcdt = await this.getAccountDcdtByIdentifier(identifier, pagination);

    return accountsDcdt.map((dcdt: any) => {
      const owner = new NftOwner();
      owner.address = dcdt.address;
      owner.balance = dcdt.balance;

      return owner;
    });
  }

  async getCollectionOwners(identifier: string, pagination: QueryPagination): Promise<CollectionAccount[] | undefined> {
    const accountsDcdt = await this.getAccountDcdtByIdentifier(identifier, pagination);

    return accountsDcdt.map((dcdt: any) => new CollectionAccount({
      address: dcdt.address,
      balance: dcdt.balance,
    }));
  }

  async getNftsInternal(pagination: QueryPagination, filter: NftFilter, identifier?: string): Promise<Nft[]> {
    if (filter.sort && filter.sort === SortCollectionNfts.rank && filter.collection) {
      const assets = await this.assetsService.getTokenAssets(filter.collection);

      filter.sort = this.getNftRankElasticKey(this.getNftRankAlgorithmFromAssets(assets));
    }

    const elasticNfts = await this.indexerService.getNfts(pagination, filter, identifier);

    const nfts: Nft[] = [];

    for (const elasticNft of elasticNfts) {
      const nft = new Nft();
      nft.identifier = elasticNft.identifier;
      nft.collection = elasticNft.token;
      nft.nonce = parseInt('0x' + nft.identifier.split('-')[2]);
      nft.timestamp = elasticNft.timestamp;

      if (elasticNft.nft_scamInfoType && elasticNft.nft_scamInfoType !== 'none') {
        nft.scamInfo = new ScamInfo({
          type: elasticNft.nft_scamInfoType,
          info: elasticNft.nft_scamInfoDescription,
        });
      }

      await this.applyExtendedAttributes(nft, elasticNft);

      const elasticNftData = elasticNft.data;
      if (elasticNftData) {
        nft.name = elasticNftData.name;
        nft.hash = TokenHelpers.getNftProof(elasticNftData.hash) ?? '';
        nft.creator = elasticNftData.creator;
        nft.royalties = elasticNftData.royalties ? elasticNftData.royalties / 100 : undefined; // 10.000 => 100%
        nft.attributes = elasticNftData.attributes;

        if (elasticNftData.uris) {
          nft.uris = elasticNftData.uris;
        }

        if (elasticNftData.tags) {
          nft.tags = elasticNftData.tags;
        }

        if (nft.uris && nft.uris.length > 0) {
          try {
            nft.url = TokenHelpers.computeNftUri(BinaryUtils.base64Decode(nft.uris[0]), this.NFT_THUMBNAIL_PREFIX);
          } catch (error) {
            this.logger.error(error);
          }
        }

        nft.isWhitelistedStorage = elasticNft.data.whiteListedStorage;
      }

      nfts.push(nft);
    }

    for (const nft of nfts) {
      const collectionProperties = await this.dcdtService.getDcdtTokenProperties(nft.collection);

      if (collectionProperties) {
        if (!nft.name) {
          nft.name = collectionProperties.name;
        }

        // @ts-ignore
        nft.type = collectionProperties.type;

        // @ts-ignore
        nft.subType = collectionProperties.subType;

        if (nft.type === NftType.MetaDCDT) {
          nft.decimals = collectionProperties.decimals;
          // @ts-ignore
          delete nft.royalties;
          // @ts-ignore
          delete nft.uris;
        }
      }
    }

    return nfts;
  }

  public isWhitelistedStorage(uris: string[] | undefined): boolean {
    if (!uris || uris.length === 0) {
      return false;
    }

    let url = '';
    try {
      url = TokenHelpers.computeNftUri(BinaryUtils.base64Decode(uris[0]), this.NFT_THUMBNAIL_PREFIX);
    } catch (error) {
      this.logger.error(`Error when computing uri from '${uris[0]}'`);
      this.logger.error(error);
      return false;
    }

    return url.startsWith(this.NFT_THUMBNAIL_PREFIX);
  }

  async getNftOwnersCount(identifier: string): Promise<number | undefined> {
    const owners = await this.cachingService.getOrSet(
      CacheInfo.NftOwnersCount(identifier).key,
      async () => await this.getNftOwnersCountRaw(identifier),
      CacheInfo.NftOwnersCount(identifier).ttl
    );

    if (owners === null) {
      return undefined;
    }

    return owners;
  }

  async getNftOwnersCountRaw(identifier: string): Promise<number | null> {
    const isNft = await this.isNft(identifier);
    if (!isNft) {
      return null;
    }

    return await this.indexerService.getNftOwnersCount(identifier);
  }

  async getNftCount(filter: NftFilter): Promise<number> {
    return await this.indexerService.getNftCount(filter);
  }

  async getNftsForAddress(address: string, queryPagination: QueryPagination, filter: NftFilter, fields?: string[], queryOptions?: NftQueryOptions, source?: DcdtDataSource): Promise<NftAccount[]> {
    let nfts = await this.dcdtAddressService.getNftsForAddress(address, filter, queryPagination, source, queryOptions);
    for (const nft of nfts) {
      await this.applyAssetsAndTicker(nft, fields);
      await this.applyPriceUsd(nft, fields);
    }

    if (queryOptions && queryOptions.withSupply) {
      const supplyNfts = nfts.filter(nft => nft.type.in(NftType.SemiFungibleDCDT, NftType.MetaDCDT));
      await this.batchApplySupply(supplyNfts, fields);
    }

    await this.batchProcessNfts(nfts, fields);

    if (this.apiConfigService.isNftExtendedAttributesEnabled() && (!fields || fields.includesSome(['score', 'rank', 'isNsfw', 'scamInfo']))) {
      const internalNfts = await this.getNftsInternalByIdentifiers(nfts.map(x => x.identifier));
      const indexedInternalNfts = internalNfts.toRecord<Nft>(x => x.identifier);
      for (const nft of nfts) {
        const indexedNft = indexedInternalNfts[nft.identifier];
        if (indexedNft) {
          nft.score = indexedNft.score;
          nft.rank = indexedNft.rank;
          nft.isNsfw = indexedNft.isNsfw;

          const scamInfo = indexedNft.scamInfo;
          if (scamInfo) {
            nft.scamInfo = new ScamInfo({
              type: scamInfo.type,
              info: scamInfo.info,
            });
          }
        }
      }
    }

    nfts = this.applyScamFilter(nfts, filter);

    for (const nft of nfts) {
      await this.applyUnlockFields(nft, fields);
    }

    return nfts;
  }

  private async getNftsInternalByIdentifiers(identifiers: string[]): Promise<Nft[]> {
    const chunks = BatchUtils.splitArrayIntoChunks(identifiers, 1024);
    const result: Nft[] = [];
    for (const identifiers of chunks) {
      const internalNfts = await this.getNftsInternal(new QueryPagination({ from: 0, size: identifiers.length }), new NftFilter({ identifiers }));

      result.push(...internalNfts);
    }

    return result;
  }

  private async applyPriceUsd(nft: NftAccount, fields?: string[]) {
    if (fields && fields.includesNone(['price', 'valueUsd'])) {
      return;
    }

    if (nft.type !== NftType.MetaDCDT) {
      return;
    }

    try {
      const prices = await this.moaTokenService.getMoaPrices();

      const price = prices[nft.collection];
      if (price) {
        nft.price = price.price;
        nft.valueUsd = price.price * NumberUtils.denominateString(nft.balance, nft.decimals);
      }
    } catch (error) {
      this.logger.error(`Unable to apply price on MetaDCDT with identifier '${nft.identifier}'`);
      this.logger.error(error);
    }
  }

  private applyScamFilter(nfts: NftAccount[], filter: NftFilter): NftAccount[] {
    if (filter.scamType) {
      nfts = nfts.filter(nft => nft.scamInfo?.type === filter.scamType);
    }

    if (filter.isScam !== undefined) {
      if (filter.isScam) {
        nfts = nfts.filter(nft => nft.scamInfo && nft.scamInfo.type);
      } else {
        nfts = nfts.filter(nft => !nft.scamInfo || !nft.scamInfo.type);
      }
    }

    return nfts;
  }

  async getNftCountForAddress(address: string, filter: NftFilter): Promise<number> {
    return await this.dcdtAddressService.getNftCountForAddressFromElastic(address, filter);
  }

  async getNftForAddress(address: string, identifier: string, fields?: string[]): Promise<NftAccount | undefined> {
    const filter = new NftFilter();
    filter.identifiers = [identifier];

    if (!TokenUtils.isNft(identifier)) {
      return undefined;
    }

    const nfts = await this.getNftsForAddress(address, new QueryPagination({ from: 0, size: 1 }), filter, fields);
    if (nfts.length === 0) {
      return undefined;
    }

    return nfts[0];
  }

  async applySupply(nft: Nft): Promise<void> {
    const { totalSupply } = await this.dcdtService.getTokenSupply(nft.identifier);

    nft.supply = totalSupply;
  }

  async getNftSupply(identifier: string): Promise<string | undefined> {
    if (identifier.split('-').length !== 3) {
      return undefined;
    }

    const nfts = await this.getNftsInternal(new QueryPagination({ from: 0, size: 1 }), new NftFilter(), identifier);
    if (nfts.length === 0) {
      return undefined;
    }

    const supply = await this.dcdtService.getTokenSupply(identifier);

    return supply.totalSupply;
  }

  async getAccountDcdtByIdentifier(identifier: string, pagination?: QueryPagination) {
    return await this.getAccountDcdtByIdentifiers([identifier], pagination);
  }

  async getAccountDcdtByIdentifiers(identifiers: string[], pagination?: QueryPagination) {
    return await this.indexerService.getAccountDcdtByIdentifiers(identifiers, pagination);
  }

  async getAccountDcdtByCollection(identifier: string, pagination?: QueryPagination) {
    return await this.indexerService.getAccountsDcdtByCollection([identifier], pagination);
  }

  // TODO: use this function to determine if a MetaDCDT is a proof if we decide to add API filters to extract all the proofs
  getNftProofHash(nft: Nft): string | undefined {
    const hashField = BinaryUtils.base64Decode(nft.hash);
    if (nft.type !== NftType.MetaDCDT || !hashField.startsWith('proof:')) {
      return undefined;
    }

    return hashField.split('proof:')[1];
  }

  private getNftRarity(elasticNft: any, algorithm: NftRankAlgorithm): NftRarity | undefined {
    const score = elasticNft[this.getNftScoreElasticKey(algorithm)];
    const rank = elasticNft[this.getNftRankElasticKey(algorithm)];

    if (!score && !rank) {
      return undefined;
    }

    return new NftRarity({ score, rank });
  }

  private async applyExtendedAttributes(nft: Nft, elasticNft: any) {
    const collectionAssets = await this.assetsService.getTokenAssets(nft.collection);
    const algorithm = this.getNftRankAlgorithmFromAssets(collectionAssets);

    nft.score = elasticNft[this.getNftScoreElasticKey(algorithm)];
    nft.rank = elasticNft[this.getNftRankElasticKey(algorithm)];

    nft.rarities = new NftRarities({
      trait: this.getNftRarity(elasticNft, NftRankAlgorithm.trait),
      statistical: this.getNftRarity(elasticNft, NftRankAlgorithm.statistical),
      jaccardDistances: this.getNftRarity(elasticNft, NftRankAlgorithm.jaccardDistances),
      openRarity: this.getNftRarity(elasticNft, NftRankAlgorithm.openRarity),
      custom: this.getNftRarity(elasticNft, NftRankAlgorithm.custom),
    });

    if (elasticNft.nft_nsfw_mark !== undefined) {
      nft.isNsfw = elasticNft.nft_nsfw_mark >= this.apiConfigService.getNftExtendedAttributesNsfwThreshold();
    }
  }

  private getNftRankAlgorithmFromAssets(assets?: TokenAssets): NftRankAlgorithm {
    return assets?.preferredRankAlgorithm ?? NftRankAlgorithm.openRarity;
  }

  private getNftRankElasticKey(algorithm: NftRankAlgorithm) {
    return `nft_rank_${algorithm}`;
  }

  private getNftScoreElasticKey(algorithm: NftRankAlgorithm) {
    return `nft_score_${algorithm}`;
  }

  private applyRedirectMedia(nft: Nft) {
    if (!nft.media || !Array.isArray(nft.media) || nft.media.length === 0) {
      return;
    }

    try {
      const network = this.apiConfigService.getNetwork();
      const defaultMediaUrl = `https://${network === 'mainnet' ? '' : `${network}-`}media.dharitri.org`;
      const defaultApiMediaUrl = `https://${network === 'mainnet' ? '' : `${network}-`}api.dharitri.org/media`;

      for (const media of nft.media) {
        if (media.url) {
          media.url = ApiUtils.replaceUri(media.url, defaultMediaUrl, this.apiConfigService.getMediaUrl());
          media.url = ApiUtils.replaceUri(media.url, defaultApiMediaUrl, this.apiConfigService.getMediaUrl());
        }
        if (media.thumbnailUrl) {
          media.thumbnailUrl = ApiUtils.replaceUri(media.thumbnailUrl, defaultMediaUrl, this.apiConfigService.getMediaUrl());
          media.thumbnailUrl = ApiUtils.replaceUri(media.thumbnailUrl, defaultApiMediaUrl, this.apiConfigService.getMediaUrl());
        }
      }
    } catch (error) {
      this.logger.error(`Error when applying redirect media for NFT with identifier '${nft.identifier}'`);
      this.logger.error(error);
    }
  }
}
