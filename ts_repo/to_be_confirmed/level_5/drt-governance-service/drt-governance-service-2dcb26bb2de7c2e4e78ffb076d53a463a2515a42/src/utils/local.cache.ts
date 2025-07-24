import LRU from 'lru-cache';
import { drtConfig } from '../config';

export default new LRU({
    max: drtConfig.localCacheMaxItems,
    // for use with tracking overall storage size
    allowStale: false,
    updateAgeOnGet: false,
    updateAgeOnHas: false,
});
