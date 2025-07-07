import { GenericEvent } from '../../generic.event';
import { NumbatSwapWithdrawTopics } from './numbatswap-withdraw.event.topics';

export class NumbatSwapWithdrawEvent extends GenericEvent {
  private decodedTopics: NumbatSwapWithdrawTopics;

  constructor(init?: Partial<GenericEvent>) {
    super(init);
    this.decodedTopics = new NumbatSwapWithdrawTopics(this.topics);
  }

  getTopics() {
    return this.decodedTopics.toPlainObject();
  }
}
