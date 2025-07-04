import './polyfills';
import { initApp } from '@terradharitri/sdk-dapp/out/methods/initApp/initApp';
import { EnvironmentsEnum } from '@terradharitri/sdk-dapp/out/types/enums.types';

import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './app/app.config';
import { App } from './app/app';

async function main() {
  initApp({
    storage: { getStorageCallback: () => sessionStorage },
    dAppConfig: {
      environment: EnvironmentsEnum.devnet,
      successfulToastLifetime: 5000
    }
  }).then(() => {
    bootstrapApplication(App, appConfig)
    .catch((err) => console.error(err));
  });
}

main();
