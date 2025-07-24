import { GenericEvent } from '../../generic.event';
import { NumbatSwapAcceptOfferTopics } from './numbatswap-acceptOffer.event.topics';

export class NumbatSwapAcceptOfferEvent extends GenericEvent {
  private decodedTopics: NumbatSwapAcceptOfferTopics;

  constructor(init?: Partial<GenericEvent>) {
    super(init);
    this.decodedTopics = new NumbatSwapAcceptOfferTopics(this.topics);
  }

  getTopics() {
    return this.decodedTopics.toPlainObject();
  }
}
