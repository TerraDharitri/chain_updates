import { RawEventType } from '../generic.types';
import { RawEvent } from '../raw.event';
import { DcdtTokenTopics } from './dcdtToken.topics';

export class DcdtLocalBurnEvent extends RawEvent {
    private decodedTopics: DcdtTokenTopics;

    constructor(init: RawEventType) {
        super(init);
        this.decodedTopics = new DcdtTokenTopics(this.topics);
    }

    getTopics() {
        return this.decodedTopics.toJSON();
    }
}
