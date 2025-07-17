import { GenericEvent } from '../../generic.event';
import { NumbatSwapBidEventsTopics } from './numbatswap-bid.event.topics';

export class NumbatSwapBidEvent extends GenericEvent {
  private decodedTopics: NumbatSwapBidEventsTopics;

  constructor(init?: Partial<GenericEvent>) {
    super(init);
    this.decodedTopics = new NumbatSwapBidEventsTopics(this.topics);
  }

  getTopics() {
    return this.decodedTopics.toPlainObject();
  }
}
