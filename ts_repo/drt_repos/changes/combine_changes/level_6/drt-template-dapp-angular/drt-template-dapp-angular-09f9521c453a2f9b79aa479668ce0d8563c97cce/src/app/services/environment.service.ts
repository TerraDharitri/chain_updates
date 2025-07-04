// import { EnvironmentsEnum } from '@terradharitri/sdk-dapp/out/types/enums.types';
import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class EnvironmentService {
  environment = 'devnet'; // Mock for now

  constructor() {}

  getEnvironment() {
    return this.environment;
  }
}
