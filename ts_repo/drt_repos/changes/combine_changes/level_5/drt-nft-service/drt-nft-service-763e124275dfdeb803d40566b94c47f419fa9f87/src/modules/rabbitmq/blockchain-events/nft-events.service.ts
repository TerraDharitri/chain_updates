import { MxApiService } from 'src/common';
import { Injectable } from '@nestjs/common';
import { NftEventEnum, NftTypeEnum } from 'src/modules/assets/models';
import { CacheEventsPublisherService } from '../cache-invalidation/cache-invalidation-publisher/change-events-publisher.service';
import { CacheEventTypeEnum, ChangedEvent } from '../cache-invalidation/events/changed.event';
import { MintEvent } from '../entities/auction/mint.event';
import { MultiTransferEvent, TransferEvent } from '../entities/auction/transfer.event';
import { FeedEventsSenderService } from './feed-events.service';
import { BurnEvent } from '../entities/auction/burn.event';
import { UpdateAttributesEvent } from '../entities/auction/update-attributes.event';

@Injectable()
export class NftEventsService {
  constructor(
    private feedEventsSenderService: FeedEventsSenderService,
    private drtApiService: MxApiService,
    private readonly cacheEventsPublisherService: CacheEventsPublisherService,
  ) { }

  public async handleNftMintEvents(mintEvents: any[], hash: string) {
    for (let event of mintEvents) {
      switch (event.identifier) {
        case NftEventEnum.DCDTNFTCreate:
          const mintEvent = new MintEvent(event);
          const createTopics = mintEvent.getTopics();
          const identifier = `${createTopics.collection}-${createTopics.nonce}`;
          const collection = await this.drtApiService.getCollectionByIdentifierForQuery(createTopics.collection, 'fields=name,type');
          if (collection?.type === NftTypeEnum.NonFungibleDCDT || collection?.type === NftTypeEnum.SemiFungibleDCDT) {
            await this.feedEventsSenderService.sendMintEvent(identifier, mintEvent, createTopics, collection);
            this.triggerCacheInvalidation(createTopics.collection, CacheEventTypeEnum.Mint);
          }
          break;

        case NftEventEnum.DCDTNFTTransfer:
          const transferEvent = new TransferEvent(event);
          const transferTopics = transferEvent.getTopics();
          const collectionInfo = await this.drtApiService.getCollectionByIdentifierForQuery(transferTopics.collection, 'fields=name,type');
          if (collectionInfo?.type === NftTypeEnum.NonFungibleDCDT || collectionInfo?.type === NftTypeEnum.SemiFungibleDCDT) {
            await this.triggerCacheInvalidationWithOwner(
              `${transferTopics.collection}-${transferTopics.nonce}`,
              CacheEventTypeEnum.OwnerChanged,
              transferEvent.getAddress(),
              transferTopics.receiverAddress.toString(),
            );
          }
          break;

        case NftEventEnum.DCDTNFTBurn:
          const burnEvent = new BurnEvent(event);
          const burnTopics = burnEvent.getTopics();
          await new Promise((resolve) => setTimeout(resolve, 500));
          const burnCollection = await this.drtApiService.getCollectionByIdentifierForQuery(burnTopics.collection, 'fields=name,type');
          if (burnCollection?.type === NftTypeEnum.NonFungibleDCDT || burnCollection?.type === NftTypeEnum.SemiFungibleDCDT) {
            await this.triggerCacheInvalidation(`${burnTopics.collection}-${burnTopics.nonce}`, CacheEventTypeEnum.AssetRefresh);
          }
          break;

        case NftEventEnum.MultiDCDTNFTTransfer:
          const multiTransferEvent = new MultiTransferEvent(event);
          multiTransferEvent.getAddress();
          const multiTransferTopics = multiTransferEvent.getTopics();
          for (const pair of multiTransferTopics.pairs) {
            if (pair.nonce !== '') {
              const collectionDetails = await this.drtApiService.getCollectionByIdentifierForQuery(pair.collection, 'fields=name,type');
              if (collectionDetails?.type === NftTypeEnum.NonFungibleDCDT || collectionDetails?.type === NftTypeEnum.SemiFungibleDCDT) {
                this.triggerCacheInvalidationWithOwner(
                  `${pair.collection}-${pair.nonce}`,
                  CacheEventTypeEnum.OwnerChanged,
                  multiTransferEvent.getAddress(),
                  multiTransferTopics.receiverAddress.toString(),
                );
              }
            }
          }
          break;

        case NftEventEnum.DCDTModifyCreator:
        case NftEventEnum.DCDTMetaDataRecreate:
        case NftEventEnum.DCDTMetaDataUpdate:
        case NftEventEnum.DCDTNFTUpdateAttributes:
          const updateEvent = new UpdateAttributesEvent(event);
          const updateTopics = updateEvent.getTopics();
          const collectionUpdateInfo = await this.drtApiService.getCollectionByIdentifierForQuery(updateTopics.collection, 'fields=name,type');
          if (collectionUpdateInfo?.type === NftTypeEnum.NonFungibleDCDT || collectionUpdateInfo?.type === NftTypeEnum.SemiFungibleDCDT) {
            await this.triggerCacheInvalidation(
              `${updateTopics.collection}-${updateTopics.nonce}`,
              CacheEventTypeEnum.AssetRefresh,
            );
          }
          break;

      }
    }
  }

  private async triggerCacheInvalidation(id: string, eventType: CacheEventTypeEnum) {
    await this.cacheEventsPublisherService.publish(
      new ChangedEvent({
        id: id,
        type: eventType,
      }),
    );
  }

  private async triggerCacheInvalidationWithOwner(id: string, eventType: CacheEventTypeEnum, address: string, receiverAddress: string) {
    await this.cacheEventsPublisherService.publish(
      new ChangedEvent({
        id: id,
        type: eventType,
        address: address,
        extraInfo: { receiverAddress: receiverAddress },
      }),
    );
  }
}
