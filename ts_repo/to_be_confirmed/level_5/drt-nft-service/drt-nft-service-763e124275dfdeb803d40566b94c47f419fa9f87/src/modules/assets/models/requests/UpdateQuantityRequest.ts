import { HandleQuantityArgs } from '../HandleQuantityArgs';

export class UpdateQuantityRequest {
  updateQuantityRoleAddress: string;
  identifier: string;
  quantity: string;
  functionName: 'DCDTNFTAddQuantity' | 'DCDTNFTBurn';
  constructor(init?: Partial<UpdateQuantityRequest>) {
    Object.assign(this, init);
  }

  static fromArgs(updateQuantityArgs: HandleQuantityArgs, functionName: 'DCDTNFTAddQuantity' | 'DCDTNFTBurn') {
    return new UpdateQuantityRequest({
      updateQuantityRoleAddress: updateQuantityArgs.addOrBurnRoleAddress,
      identifier: updateQuantityArgs.identifier,
      quantity: updateQuantityArgs.quantity,
      functionName: functionName,
    });
  }
}
