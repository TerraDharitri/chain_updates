import * as crypto from 'crypto-js';
import { Nft } from "src/endpoints/nfts/entities/nft";
import { NftType } from "src/endpoints/nfts/entities/nft.type";
import { CollectionRoles } from "src/endpoints/tokens/entities/collection.roles";
import { TokenRoles } from "src/endpoints/tokens/entities/token.roles";
import { ApiUtils } from '@terradharitri/sdk-nestjs-http';
import '@terradharitri/sdk-nestjs-common/lib/utils/extensions/string.extensions';
import { BinaryUtils } from "@terradharitri/sdk-nestjs-common";

export class TokenHelpers {
  static canBool(string: string) {
    return string.split('-').pop() === 'true';
  }

  static computeNftUri(uri: string, prefix: string) {
    uri = ApiUtils.replaceUri(uri, 'https://ipfs.io/ipfs', prefix);
    uri = ApiUtils.replaceUri(uri, 'https://gateway.ipfs.io/ipfs', prefix);
    uri = ApiUtils.replaceUri(uri, 'https://gateway.pinata.cloud/ipfs', prefix);
    uri = ApiUtils.replaceUri(uri, 'https://dweb.link/ipfs', prefix);
    uri = ApiUtils.replaceUri(uri, 'ipfs:/', prefix);
    uri = ApiUtils.replaceUri(uri, 'https://media.dharitri.org/nfts/asset', prefix);
    uri = ApiUtils.replaceUri(uri, 'https://devnet-media.dharitri.org/nfts/asset', prefix);
    uri = ApiUtils.replaceUri(uri, 'https://testnet-media.dharitri.org/nfts/asset', prefix);
    uri = ApiUtils.replaceUri(uri, 'https://api.dharitri.org/media/nfts/asset', prefix);
    uri = ApiUtils.replaceUri(uri, 'https://devnet-api.dharitri.org/media/nfts/asset', prefix);
    uri = ApiUtils.replaceUri(uri, 'https://testnet-api.dharitri.org/media/nfts/asset', prefix);
    uri = uri.replace(/https\:\/\/\w*\.mypinata\.cloud\/ipfs/, prefix);

    if (uri.endsWith('.ipfs.dweb.link')) {
      const id = uri.removeSuffix('.ipfs.dweb.link').removePrefix('https://');
      uri = `${prefix}/${id}`;
    }

    return uri;
  }

  static getUrlHash(url: string) {
    return crypto.SHA256(url.trim()).toString().slice(0, 8);
  }

  static getThumbnailUrlIdentifier(nftIdentifier: string, fileUrl: string) {
    const collectionIdentifier = nftIdentifier.split('-').slice(0, 2).join('-');
    const urlHash = TokenHelpers.getUrlHash(fileUrl);

    return `${collectionIdentifier}-${urlHash}`;
  }

  static needsDefaultMedia(nft: Nft): boolean {
    if (nft.type === NftType.MetaDCDT) {
      return false;
    }

    if (nft.media && nft.media.length > 0) {
      return false;
    }

    return true;
  }

  static setTokenRole(tokenRoles: TokenRoles, role: string) {
    tokenRoles.roles.push(role);

    switch (role) {
      case 'DCDTRoleLocalMint':
        tokenRoles.canLocalMint = true;
        break;
      case 'DCDTRoleLocalBurn':
        tokenRoles.canLocalBurn = true;
        break;
    }
  }

  static setCollectionRole(tokenRoles: CollectionRoles, role: string) {
    tokenRoles.roles.push(role);

    switch (role) {
      case 'DCDTRoleNFTCreate':
        tokenRoles.canCreate = true;
        break;
      case 'DCDTRoleNFTBurn':
        tokenRoles.canBurn = true;
        break;
      case 'DCDTRoleNFTAddQuantity':
        tokenRoles.canAddQuantity = true;
        break;
      case 'DCDTRoleNFTAddURI':
        tokenRoles.canAddUri = true;
        break;
      case 'DCDTRoleNFTUpdateAttributes':
        tokenRoles.canUpdateAttributes = true;
        break;
      case 'DCDTTransferRole':
        tokenRoles.canTransfer = true;
        break;
    }
  }

  static tokenNonce(tokenID: string): number {
    const tokenNonceHex = tokenID.split('-')[2];
    return parseInt(tokenNonceHex, 16);
  }

  static getCollectionIdentifier(nftIdentifier: string): string {
    return nftIdentifier.split('-').slice(0, 2).join('-');
  }

  static getNftProof(hash: string): string | undefined {
    if (!hash) {
      return undefined;
    }

    const decodedHex = BinaryUtils.base64Decode(hash);
    if (decodedHex.startsWith('proof:')) {
      return decodedHex;
    } else {
      return hash;
    }
  }
}
