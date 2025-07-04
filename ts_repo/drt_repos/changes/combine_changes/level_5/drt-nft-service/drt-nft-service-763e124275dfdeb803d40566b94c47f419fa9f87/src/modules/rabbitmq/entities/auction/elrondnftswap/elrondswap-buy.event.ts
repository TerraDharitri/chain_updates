import { GenericEvent } from '../../generic.event';
import { NumbatSwapBuyTopics } from './numbatswap-buy.event.topics';

export class NumbatSwapBuyEvent extends GenericEvent {
  private decodedTopics: NumbatSwapBuyTopics;

  constructor(init?: Partial<GenericEvent>) {
    super(init);
    this.decodedTopics = new NumbatSwapBuyTopics(this.topics);
  }

  getTopics() {
    return this.decodedTopics.toPlainObject();
  }
}
