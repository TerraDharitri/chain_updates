export class DcdtAddressRoles {
  constructor(init?: Partial<DcdtAddressRoles>) {
    Object.assign(this, init);
  }
  roles!: { [key: string]: any };
}
