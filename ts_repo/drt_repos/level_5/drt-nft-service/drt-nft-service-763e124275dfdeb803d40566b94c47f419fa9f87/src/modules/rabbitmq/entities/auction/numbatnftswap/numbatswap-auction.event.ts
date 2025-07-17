import { GenericEvent } from '../../generic.event';
import { NumbatSwapAuctionTopics } from './numbatswap-auction.event.topics';

export class NumbatSwapAuctionEvent extends GenericEvent {
  private decodedTopics: NumbatSwapAuctionTopics;

  constructor(init?: Partial<GenericEvent>) {
    super(init);
    this.decodedTopics = new NumbatSwapAuctionTopics(this.topics);
  }

  getTopics() {
    return this.decodedTopics.toPlainObject();
  }
}
