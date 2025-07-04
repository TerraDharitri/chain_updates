import { GenericEvent } from '../../generic.event';
import { NumbatSwapUpdateTopics } from './numbatswap-updateAuction.event.topics';

export class NumbatSwapUpdateEvent extends GenericEvent {
  private decodedTopics: NumbatSwapUpdateTopics;

  constructor(init?: Partial<GenericEvent>) {
    super(init);
    this.decodedTopics = new NumbatSwapUpdateTopics(this.topics);
  }

  getTopics() {
    return this.decodedTopics.toPlainObject();
  }
}
