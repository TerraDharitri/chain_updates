import { Injectable } from "@nestjs/common";
import { CacheInfo } from "src/utils/cache.info";
import { TokenAssets } from "src/common/assets/entities/token.assets";
import { AccountAssets } from "./entities/account.assets";
import { CacheService } from "@terradharitri/sdk-nestjs-cache";
import { MoaPair } from "src/endpoints/moa/entities/moa.pair";
import { Identity } from "src/endpoints/identities/entities/identity";
import { MoaFarm } from "src/endpoints/moa/entities/moa.farm";
import { MoaSettings } from "src/endpoints/moa/entities/moa.settings";
import { DnsContracts } from "src/utils/dns.contracts";
import { NftRank } from "./entities/nft.rank";
import { MoaStakingProxy } from "src/endpoints/moa/entities/moa.staking.proxy";
import { Provider } from "src/endpoints/providers/entities/provider";
import { ApiService } from "@terradharitri/sdk-nestjs-http";
import { ApiConfigService } from "../api-config/api.config.service";
import { KeybaseIdentity } from "../keybase/entities/keybase.identity";

@Injectable()
export class AssetsService {
  constructor(
    private readonly apiConfigService: ApiConfigService,
    private readonly apiService: ApiService,
    private readonly cachingService: CacheService,
  ) { }

  async getAllTokenAssets(): Promise<{ [key: string]: TokenAssets }> {
    return await this.cachingService.getOrSet(
      CacheInfo.TokenAssets.key,
      async () => await this.getAllTokenAssetsRaw(),
      CacheInfo.TokenAssets.ttl,
    );
  }

  async getAllTokenAssetsRaw(): Promise<{ [key: string]: TokenAssets }> {
    if (!this.apiConfigService.isAssetsCdnFeatureEnabled()) {
      return {};
    }

    const assetsCdnUrl = this.apiConfigService.getAssetsCdnUrl();
    const network = this.apiConfigService.getNetwork();

    const { data: assetsRaw } = await this.apiService.get(`${assetsCdnUrl}/${network}/tokens`);

    const assets: { [key: string]: TokenAssets } = {};
    for (const asset of assetsRaw) {
      const { identifier, ...details } = asset;
      assets[identifier] = new TokenAssets(details);
    }

    return assets;
  }

  async getCollectionRanks(identifier: string): Promise<NftRank[] | undefined> {
    const allCollectionRanks = await this.getAllCollectionRanks();
    return allCollectionRanks[identifier];
  }

  async getAllCollectionRanks(): Promise<{ [key: string]: NftRank[] }> {
    return await this.cachingService.getOrSet(
      CacheInfo.CollectionRanks.key,
      async () => await this.getAllCollectionRanksRaw(),
      CacheInfo.CollectionRanks.ttl,
    );
  }

  async getAllCollectionRanksRaw(): Promise<{ [key: string]: NftRank[] }> {
    if (!this.apiConfigService.isAssetsCdnFeatureEnabled()) {
      return {};
    }

    const assetsCdnUrl = this.apiConfigService.getAssetsCdnUrl();
    const network = this.apiConfigService.getNetwork();

    const { data: assets } = await this.apiService.get(`${assetsCdnUrl}/${network}/tokens`);

    const result: { [key: string]: NftRank[] } = {};

    for (const asset of assets) {
      if (asset.ranks && asset.ranks.length > 0) {
        result[asset.identifier] = asset.ranks.map((rank: any) => new NftRank({
          identifier: rank.identifier,
          rank: rank.rank,
        }));
      }
    }

    return result;
  }

  async getAllAccountAssets(): Promise<{ [key: string]: AccountAssets }> {
    return await this.cachingService.getOrSet(
      CacheInfo.AccountAssets.key,
      async () => await this.getAllAccountAssetsRaw(),
      CacheInfo.AccountAssets.ttl,
    );
  }

  async getAllAccountAssetsRaw(providers?: Provider[], identities?: Identity[], pairs?: MoaPair[], farms?: MoaFarm[], moaSettings?: MoaSettings, stakingProxies?: MoaStakingProxy[]): Promise<{ [key: string]: AccountAssets }> {
    if (!this.apiConfigService.isAssetsCdnFeatureEnabled()) {
      return {};
    }

    const assetsCdnUrl = this.apiConfigService.getAssetsCdnUrl();
    const network = this.apiConfigService.getNetwork();

    const { data: assets } = await this.apiService.get(`${assetsCdnUrl}/${network}/accounts`);

    const allAssets: { [key: string]: AccountAssets } = {};
    for (const asset of assets) {
      const { address, ...details } = asset;
      allAssets[address] = new AccountAssets(details);
    }

    // Populate additional assets from other sources if available
    if (providers && identities) {
      for (const provider of providers) {
        const identity = identities.find(x => x.identity === provider.identity);
        if (!identity) {
          continue;
        }

        allAssets[provider.provider] = new AccountAssets({
          name: `Staking: ${identity.name ?? ''}`,
          description: identity.description ?? '',
          iconPng: identity.avatar,
          tags: ['staking', 'provider'],
        });
      }
    }

    if (pairs) {
      for (const pair of pairs) {
        allAssets[pair.address] = this.createAccountAsset(
          `xExchange: ${pair.baseSymbol}/${pair.quoteSymbol} Liquidity Pool`,
          ['dharitrix', 'liquiditypool']
        );
      }
    }

    if (farms) {
      for (const farm of farms) {
        allAssets[farm.address] = this.createAccountAsset(
          `xExchange: ${farm.name} Farm`,
          ['dharitrix', 'farm']
        );
      }
    }

    if (moaSettings) {
      for (const [index, wrapContract] of moaSettings.wrapContracts.entries()) {
        allAssets[wrapContract] = this.createAccountAsset(
          `DCDT: WrappedREWA Contract Shard ${index}`,
          ['dharitrix', 'wrewa']
        );
      }

      allAssets[moaSettings.lockedAssetContract] = this.createAccountAsset(
        `xExchange: Locked asset Contract`,
        ['dharitrix', 'lockedasset']
      );

      allAssets[moaSettings.distributionContract] = this.createAccountAsset(
        `xExchange: Distribution Contract`,
        ['dharitrix', 'lockedasset']
      );
    }

    if (stakingProxies) {
      for (const stakingProxy of stakingProxies) {
        allAssets[stakingProxy.address] = this.createAccountAsset(
          `xExchange: ${stakingProxy.dualYieldTokenName} Contract`,
          ['dharitrix', 'metastaking']
        );
      }
    }

    for (const [index, address] of DnsContracts.addresses.entries()) {
      allAssets[address] = new AccountAssets({
        name: `Dharitri DNS: Contract ${index}`,
        tags: ['dns'],
        icon: 'dharitri',
      });
    }

    return allAssets;
  }

  async getTokenAssets(tokenIdentifier: string): Promise<TokenAssets | undefined> {
    const assets = await this.getAllTokenAssets();
    return assets[tokenIdentifier];
  }

  async getAllIdentitiesRaw(): Promise<{ [key: string]: KeybaseIdentity }> {
    if (!this.apiConfigService.isAssetsCdnFeatureEnabled()) {
      return {};
    }

    const assetsCdnUrl = this.apiConfigService.getAssetsCdnUrl();
    const network = this.apiConfigService.getNetwork();

    const { data: assets } = await this.apiService.get(`${assetsCdnUrl}/${network}/identities`);

    const allAssets: { [key: string]: KeybaseIdentity } = {};
    for (const asset of assets) {
      allAssets[asset.identity] = new KeybaseIdentity(asset);
    }

    return allAssets;
  }

  async getIdentityInfo(identity: string): Promise<KeybaseIdentity | null> {
    const allIdentities = await this.getAllIdentitiesRaw();
    return allIdentities[identity] || null;
  }

  createAccountAsset(name: string, tags: string[]): AccountAssets {
    return new AccountAssets({
      name: name,
      tags: tags,
      iconSvg: 'https://raw.githubusercontent.com/TerraDharitri/drt-assets/master/accounts/icons/dharitrix.svg',
      iconPng: 'https://raw.githubusercontent.com/TerraDharitri/drt-assets/master/accounts/icons/dharitrix.png',
    });
  }
}
