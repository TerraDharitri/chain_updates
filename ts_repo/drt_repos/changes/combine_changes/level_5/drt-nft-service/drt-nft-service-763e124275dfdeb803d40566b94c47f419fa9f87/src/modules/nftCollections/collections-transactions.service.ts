import {
  Address,
  SemiFungibleSpecialRoleInput,
  SpecialRoleInput,
  TokenManagementTransactionsFactory,
  TransactionsFactoryConfig,
} from '@terradharitri/sdk-core';
import { Injectable } from '@nestjs/common';
import { MxApiService } from 'src/common';
import { drtConfig } from 'src/config';
import { NftTypeEnum } from '../assets/models';
import { TransactionNode } from '../common/transaction';
import { IssueCollectionRequest, SetNftRolesRequest } from './models/requests';

@Injectable()
export class CollectionsTransactionsService {
  constructor(private apiService: MxApiService) {}
  async issueToken(ownerAddress: string, request: IssueCollectionRequest) {
    const factory = new TokenManagementTransactionsFactory({ config: new TransactionsFactoryConfig({ chainID: drtConfig.chainID }) });
    if (request.collectionType === 'issueNonFungible') {
      const transaction = factory.createTransactionForIssuingNonFungible(Address.newFromBech32(ownerAddress), {
        tokenName: request.tokenName,
        tokenTicker: request.tokenTicker,
        canFreeze: request.canFreeze,
        canWipe: request.canWipe,
        canPause: request.canPause,
        canTransferNFTCreateRole: request.canTransferNFTCreateRole,
        canChangeOwner: request.canChangeOwner,
        canUpgrade: request.canUpgrade,
        canAddSpecialRoles: request.canAddSpecialRoles,
      });
      return transaction.toPlainObject();
    }
    const transaction = factory.createTransactionForIssuingSemiFungible(Address.newFromBech32(ownerAddress), {
      tokenName: request.tokenName,
      tokenTicker: request.tokenTicker,
      canFreeze: request.canFreeze,
      canWipe: request.canWipe,
      canPause: request.canPause,
      canTransferNFTCreateRole: request.canTransferNFTCreateRole,
      canChangeOwner: request.canChangeOwner,
      canUpgrade: request.canUpgrade,
      canAddSpecialRoles: request.canAddSpecialRoles,
    });
    return transaction.toPlainObject();
  }

  async setNftRoles(ownerAddress: string, args: SetNftRolesRequest): Promise<TransactionNode> {
    const factory = new TokenManagementTransactionsFactory({
      config: new TransactionsFactoryConfig({ chainID: drtConfig.chainID }),
    });

    const collection = await this.apiService.getCollectionForIdentifier(args.collection);
    const userAddress = Address.newFromBech32(ownerAddress);

    if (collection.type === NftTypeEnum.NonFungibleDCDT) {
      const nftInput: SpecialRoleInput = {
        user: userAddress,
        tokenIdentifier: args.collection,
        addRoleNFTCreate: args.roles.includes('DCDTRoleNFTCreate'),
        addRoleNFTBurn: args.roles.includes('DCDTRoleNFTBurn'),
        addRoleNFTUpdateAttributes: args.roles.includes('DCDTRoleNFTUpdateAttributes'),
        addRoleNFTAddURI: args.roles.includes('DCDTRoleNFTAddURI'),
        addRoleDCDTTransferRole: args.roles.includes('DCDTTransferRole'),
        addRoleDCDTModifyCreator: args.roles.includes('DCDTRoleModifyCreator'),
        addRoleNFTRecreate: args.roles.includes('DCDTRoleNFTRecreate'),
        addRoleDCDTModifyRoyalties: args.roles.includes('DCDTRoleModifyRoyalties'),
      };

      return factory.createTransactionForSettingSpecialRoleOnNonFungibleToken(userAddress, nftInput).toPlainObject();
    }

    const sftInput: SemiFungibleSpecialRoleInput = {
      user: userAddress,
      tokenIdentifier: args.collection,
      addRoleNFTCreate: args.roles.includes('DCDTRoleNFTCreate'),
      addRoleNFTBurn: args.roles.includes('DCDTRoleNFTBurn'),
      addRoleDCDTTransferRole: args.roles.includes('DCDTTransferRole'),
      addRoleNFTAddQuantity: args.roles.includes('DCDTRoleNFTAddQuantity'),
    };

    return factory.createTransactionForSettingSpecialRoleOnSemiFungibleToken(userAddress, sftInput).toPlainObject();
  }
}
