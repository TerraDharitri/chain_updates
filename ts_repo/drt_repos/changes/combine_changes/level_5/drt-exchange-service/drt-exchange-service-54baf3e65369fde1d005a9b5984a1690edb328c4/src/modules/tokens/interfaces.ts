export interface ITokenComputeService {
    getRewaPriceInUSD(): Promise<string>;
    computeTokenPriceDerivedREWA(
        tokenID: string,
        pairsNotToVisit: [],
    ): Promise<string>;
    computeTokenPriceDerivedUSD(tokenID: string): Promise<string>;
}
