import GenericGuardianProvider from "./genericGuardianProvider";
import TCSGuardianProvider from "./providers/TCSGuardianProvider";

export const DEFAULT_SERVICE_ID = "MultiversXTCSService";

enum NetworksEnum {
  testnet = "testnet",
  devnet = "devnet",
  mainnet = "mainnet",
}

type ServiceNetworkType = Record<NetworksEnum | string, string>;

interface IInjectableNetworkUrl {
  /**
   * Network identifier.
   */
  networkId: string;
  /**
   * Provider service network url.
   */
  url: string;
}

interface IProviderInfo {
  /**
   * `serviceId` is the unique identifier of the `IProviderInfo` object.
   */
  serviceId: string;
  provider: typeof GenericGuardianProvider;
  providerServiceNetworkUrls: ServiceNetworkType;
}

class GuardianProvidersResolver {
  protected static providers: Array<IProviderInfo> = [
    {
      serviceId: DEFAULT_SERVICE_ID,
      provider: TCSGuardianProvider,
      providerServiceNetworkUrls: {
        testnet: "https://testnet-tcs-api.multiversx.com",
        devnet: "https://devnet-tools.multiversx.com",
        'devnet-old': "https://devnet-old-tools.multiversx.com",
        mainnet: "https://tools.multiversx.com",
        battlenet: "https://tools.multiversx.com"
      },
    },
  ];

  static getProviderByServiceId(serviceId: string): IProviderInfo | undefined {
    const result = this.providers.find(
      (provider) => provider.serviceId === serviceId
    );
    return result;
  }

  /**
   * Extends or overrides the service networks url list based on network.
   */
  static extendProviderInfoNetworkUrls({
    serviceId,
    network,
  }: {
    serviceId: IProviderInfo["serviceId"];
    network: IInjectableNetworkUrl;
  }) {
    this.providers = this.providers.map((provider) => {
      if (provider.serviceId === serviceId) {
        provider.providerServiceNetworkUrls = {
          ...provider.providerServiceNetworkUrls,
          [network.networkId]: network.url,
        };
      }
      return provider;
    });
  }

  public static get defaultServiceId(): string {
    return DEFAULT_SERVICE_ID;
  }
}

export default GuardianProvidersResolver;
