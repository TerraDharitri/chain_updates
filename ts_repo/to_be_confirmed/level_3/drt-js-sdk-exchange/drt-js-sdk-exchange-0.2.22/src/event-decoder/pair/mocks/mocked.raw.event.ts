import { RawEventType } from '../../generic.types';

export const rawAddLiquidityEvent: RawEventType = {
    address: 'drt1qqqqqqqqqqqqqpgquu5rsa4ee6l4azz6vdu4hjp8z4p6tt8m0n4sptujw6',
    identifier: 'addLiquidity',
    topics: [
        'YWRkX2xpcXVpZGl0eQ==',
        'V1JFV0EtZDdjNmJi',
        'TU9BLWRjMjg5Yw==',
        'NEq8RBGc/KziU94F4zwBeWwS+W87zFK1BLm8K5aSfOs=',
        'AaY=',
    ],
    data: 'NEq8RBGc/KziU94F4zwBeWwS+W87zFK1BLm8K5aSfOsAAAAMV1JFV0EtZDdjNmJiAAAACA3gtrOnZAAAAAAACk1PQS1kYzI4OWMAAAALAd2dBVqPm1WQus0AAAAOUkVXQU1PQS1jMjliMGUAAAAIGVJP+UphDWMAAAAJVliatA0RsWe1AAAACS9StZK8KUai8wAAAAwGXKfVif5m3k+kyuMAAAAAAAe2mQAAAAAAAAGmAAAAAGJxVOY=',
};

export const rawRemoveLiquidityEvent: RawEventType = {
    address: 'drt1qqqqqqqqqqqqqpgquu5rsa4ee6l4azz6vdu4hjp8z4p6tt8m0n4sptujw6',
    identifier: 'removeLiquidity',
    topics: [
        'cmVtb3ZlX2xpcXVpZGl0eQ==',
        'V1JFV0EtZDdjNmJi',
        'TU9BLWRjMjg5Yw==',
        'NEq8RBGc/KziU94F4zwBeWwS+W87zFK1BLm8K5aSfOs=',
        'AbA=',
    ],
    data: 'NEq8RBGc/KziU94F4zwBeWwS+W87zFK1BLm8K5aSfOsAAAAMV1JFV0EtZDdjNmJiAAAACA53PRJHvS1gAAAACk1PQS1kYzI4OWMAAAALAiagMnZ07SEyyK0AAAAOUkVXQU1PQS1jMjliMGUAAAAIG8FtZ07IAAAAAAAJVjzZRqXC6We1AAAACSzyHHmht78awwAAAAwGrtEBQrHe0Ec3BkEAAAAAAAfnLwAAAAAAAAGwAAAAAGJyeX4=',
};

export const rawSwapFixedInputEvent: RawEventType = {
    address: 'drt1qqqqqqqqqqqqqpgquu5rsa4ee6l4azz6vdu4hjp8z4p6tt8m0n4sptujw6',
    identifier: 'swapTokensFixedInput',
    topics: [
        'c3dhcA==',
        'V1JFV0EtZDdjNmJi',
        'TU9BLWRjMjg5Yw==',
        'NEq8RBGc/KziU94F4zwBeWwS+W87zFK1BLm8K5aSfOs=',
        'AbA=',
    ],
    data: 'NEq8RBGc/KziU94F4zwBeWwS+W87zFK1BLm8K5aSfOsAAAAMV1JFV0EtZDdjNmJiAAAACA3gtrOnZAAAAAAACk1PQS1kYzI4OWMAAAALAg4G/P0nHBAQb8AAAAAHAca/UmNAAAAAAAks//0wVV8jGsMAAAAMBqzCttHN8WBD5crgAAAAAAAH5zwAAAAAAAABsAAAAABicnnM',
};

export const rawSwapFixedOutputEvent: RawEventType = {
    address: 'drt1qqqqqqqqqqqqqpgquu5rsa4ee6l4azz6vdu4hjp8z4p6tt8m0n4sptujw6',
    identifier: 'swapTokensFixedOutput',
    topics: [
        'c3dhcA==',
        'V1JFV0EtZDdjNmJi',
        'TU9BLWRjMjg5Yw==',
        'NEq8RBGc/KziU94F4zwBeWwS+W87zFK1BLm8K5aSfOs=',
        'AbA=',
    ],
    data: 'NEq8RBGc/KziU94F4zwBeWwS+W87zFK1BLm8K5aSfOsAAAAMV1JFV0EtZDdjNmJiAAAACA3gtrOl8NT2AAAACk1PQS1kYzI4OWMAAAALAgzDVWkctbXugAAAAAAHAca/UmMQfQAAAAktDd3nCQUT77kAAAAMBqq1sDH4RJFY9piUAAAAAAAH50wAAAAAAAABsAAAAABicnoy',
};

export const rawSwapNoFeeEvent: RawEventType = {
    address: 'drt1qqqqqqqqqqqqqpgquu5rsa4ee6l4azz6vdu4hjp8z4p6tt8m0n4sptujw6',
    identifier: 'swapNoFeeAndForward',
    topics: [
        'c3dhcF9ub19mZWVfYW5kX2ZvcndhcmQ=',
        'TU9BLWRjMjg5Yw==',
        'AAAAAAAAAAAFAAa9xh677HGbB7Sn6/0fshXAcG48fOs=',
        'AbA=',
    ],
    data: 'AAAAAAAAAAAFAAa9xh677HGbB7Sn6/0fshXAcG48fOsAAAAMV1JFV0EtZDdjNmJiAAAAB5pZwTC9lHIAAAAKTU9BLWRjMjg5YwAAAAoW1rGIh5yGJp+9AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAfnWgAAAAAAAAGwAAAAAGJyeoY=',
};
