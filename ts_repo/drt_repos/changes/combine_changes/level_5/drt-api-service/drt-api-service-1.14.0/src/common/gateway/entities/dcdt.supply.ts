
export class DcdtSupply {
  constructor(init?: Partial<DcdtSupply>) {
    Object.assign(this, init);
  }

  supply: string = '0';
  minted: string = '0';
  burned: string = '0';
  initialMinted: string = '0';
}
