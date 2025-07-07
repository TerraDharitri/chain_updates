export class DataApiToken {
  constructor(init?: Partial<DataApiToken>) {
    Object.assign(this, init);
  }

  identifier: string = '';
  market: 'cex' | 'dharitrix' | 'hatom' | 'xoxno' = 'cex';
}
