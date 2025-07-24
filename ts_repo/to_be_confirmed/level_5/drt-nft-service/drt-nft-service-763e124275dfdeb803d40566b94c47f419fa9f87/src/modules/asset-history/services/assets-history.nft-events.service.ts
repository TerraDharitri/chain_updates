import { Injectable } from '@nestjs/common';
import { drtConfig } from 'src/config';
import { NftEventEnum, AssetActionEnum } from 'src/modules/assets/models';
import { AssetHistoryInput as AssetHistoryLogInput } from '../models/asset-history-log-input';

@Injectable()
export class AssetsHistoryNftEventService {
  constructor() {}

  mapNftEventLog(nonce: string, eventType: string, mainEvent: any): AssetHistoryLogInput {
    const event = mainEvent.events.find((event) => event.identifier === eventType);
    const encodedNonce = Buffer.from(nonce, 'hex').toString('base64');
    const transferEvent = mainEvent.events.find(
      (event) =>
        (event.identifier === NftEventEnum.DCDTNFTTransfer || event.identifier === NftEventEnum.MultiDCDTNFTTransfer) &&
        event.topics[1] === encodedNonce,
    );

    switch (eventType) {
      case NftEventEnum.DCDTNFTAddQuantity: {
        return new AssetHistoryLogInput({
          event: mainEvent,
          action: AssetActionEnum.Added,
          address: mainEvent.address,
          itemsCount: mainEvent.events[0].topics[2],
        });
      }
      case NftEventEnum.DCDTNFTCreate: {
        return new AssetHistoryLogInput({
          event: mainEvent,
          action: AssetActionEnum.Created,
          address: event.address,
          itemsCount: event.topics[2],
          sender: transferEvent?.topics[3]?.base64ToBech32(),
        });
      }
      case NftEventEnum.DCDTNFTTransfer: {
        if (
          mainEvent.address === mainEvent?.events[0].address &&
          transferEvent.topics[3].base64ToBech32() !== drtConfig.nftMarketplaceAddress
        ) {
          return new AssetHistoryLogInput({
            event: mainEvent,
            action: AssetActionEnum.Received,
            address: transferEvent.topics[3].base64ToBech32(),
            itemsCount: transferEvent.topics[2],
            sender: transferEvent.address,
          });
        }
      }
      case NftEventEnum.MultiDCDTNFTTransfer: {
        const senderAddress = transferEvent.address;
        const receiverAddress = transferEvent.topics[3].base64ToBech32();
        if (senderAddress !== receiverAddress) {
          return new AssetHistoryLogInput({
            event: mainEvent,
            action: AssetActionEnum.Received,
            address: receiverAddress,
            itemsCount: transferEvent.topics[2],
            sender: senderAddress,
          });
        }
      }
      case NftEventEnum.DCDTNFTBurn:
      case NftEventEnum.DCDTNFTUpdateAttributes: {
        break;
      }
      default: {
        return this.mapNftEventLog(nonce, transferEvent.identifier, mainEvent);
      }
    }
  }
}
