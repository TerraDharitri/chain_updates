export class DrtApiAbout {
  appVersion: string;
  pluginsVersion: string;
  network: string;
  cluster: string;
  version: string;
  scamEngineVersion: string = '1.0.0';

  constructor(init?: Partial<DrtApiAbout>) {
    Object.assign(this, init);
  }
}
