export interface EntrypointConfig {
    networkProviderUrl: string;
    networkProviderKind: string;
    chainId: string;
}

export class TestnetEntrypointConfig {
    networkProviderUrl: string;
    networkProviderKind: string;
    chainId: string;

    constructor({
        networkProviderUrl = "https://testnet-api.dharitri.org",
        networkProviderKind = "api",
        chainId = "T",
    }: Partial<EntrypointConfig> = {}) {
        this.networkProviderUrl = networkProviderUrl;
        this.networkProviderKind = networkProviderKind;
        this.chainId = chainId;
    }
}

export class DevnetEntrypointConfig {
    networkProviderUrl: string;
    networkProviderKind: string;
    chainId: string;
    constructor({
        networkProviderUrl = "https://devnet-api.dharitri.org",
        networkProviderKind = "api",
        chainId = "D",
    }: Partial<EntrypointConfig> = {}) {
        this.networkProviderUrl = networkProviderUrl;
        this.networkProviderKind = networkProviderKind;
        this.chainId = chainId;
    }
}

export class MainnetEntrypointConfig {
    networkProviderUrl: string;
    networkProviderKind: string;
    chainId: string;

    constructor({
        networkProviderUrl = "https://api.dharitri.org",
        networkProviderKind = "api",
        chainId = "1",
    }: Partial<EntrypointConfig> = {}) {
        this.networkProviderUrl = networkProviderUrl;
        this.networkProviderKind = networkProviderKind;
        this.chainId = chainId;
    }
}

export class LocalnetEntrypointConfig {
    networkProviderUrl: string;
    networkProviderKind: string;
    chainId: string;

    constructor({
        networkProviderUrl = "http://localhost:7950",
        networkProviderKind = "proxy",
        chainId = "localnet",
    }: Partial<EntrypointConfig> = {}) {
        this.networkProviderUrl = networkProviderUrl;
        this.networkProviderKind = networkProviderKind;
        this.chainId = chainId;
    }
}
