declare global {
    namespace NodeJS {
      interface ProcessEnv {
        DRT_WALLET: string;
        DRT_NETWORK: string;
      }
    }
}

export { };
