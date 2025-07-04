import { RawEventType } from '../../generic.types';

export const rawStakeEvent: RawEventType = {
    identifier: 'stakeFarm',
    address: 'drt1qqqqqqqqqqqqqpgqp2r407yzwvreva2k5wp8ul8vks7ukg8j295qc0pcj0',
    data: 'AAAAClVUSy1iOTc0ODAAAAAIDeC2s6dkAAAAAAALU1VUSy0xNTk4OGEAAAAAAAAAAQAAAAgN4Lazp2QAAAAAAAgN4Lazp2QAAAAAAApVVEstYjk3NDgwAAAAAAAAABQAAAAAAAAAAAAAAAgN4Lazp2QAAAA=',
    topics: [
        'ZW50ZXJfZmFybQ==',
        'YLtgEaeB64tT1h15BHRQrapGl34gI5MSlRJiwQV9EJA=',
        'DVg=',
        'Poy5',
        'ZSQdOA==',
        'VVRLLWI5NzQ4MA==',
    ],
};

export const rawUnstakeEvent: RawEventType = {
    identifier: 'unstakeFarm',
    address: 'drt1qqqqqqqqqqqqqpgqp2r407yzwvreva2k5wp8ul8vks7ukg8j295qc0pcj0',
    data: 'AAAAC1NVVEstMTU5ODhhAAAACA3gtrOnZAAAAAAAC1NVVEstMTU5ODhhAAAAAAAAAAQAAAAIDeC2s6dkAAAAAAAAAAAAClVUSy1iOTc0ODAAAAAAAAAAAAAAAAYbQPIQIEAAAAADHpekAAAAOAAAAARuJ3AxAAAAAAAAAAgN4Lazp2QAAGC7YBGngeuLU9YdeQR0UK2qRpd+ICOTEpUSYsEFfRCQ',
    topics: [
        'ZXhpdF9mYXJt',
        'YLtgEaeB64tT1h15BHRQrapGl34gI5MSlRJiwQV9EJA=',
        'DXk=',
        'Pyb1',
        'ZSe7og==',
        'U1VUSy0xNTk4OGE=',
    ],
};

export const rawStakeClaimRewardsEvent: RawEventType = {
    identifier: 'claimRewards',
    address: 'drt1qqqqqqqqqqqqqpgqp2r407yzwvreva2k5wp8ul8vks7ukg8j295qc0pcj0',
    data: 'AAAAC1NVVEstMTU5ODhhAAAAAAAAAAMAAAAIDeC2s6dkAAAAAAALU1VUSy0xNTk4OGEAAAAAAAAABAAAAAgN4Lazp2QAAAAAAAgN4Lazp2QAAAAAAApVVEstYjk3NDgwAAAAAAAAAAAAAAAHAhqKqTQwAAAAAAMYHmoAAAA4AAAABErcNXEAAAAAAAAACA3gtrOnZAAAYLtgEaeB64tT1h15BHRQrapGl34gI5MSlRJiwQV9EJAAAAA4AAAABG4ncDEAAAAAAAAACA3gtrOnZAAAYLtgEaeB64tT1h15BHRQrapGl34gI5MSlRJiwQV9EJAA',
    topics: [
        'Y2xhaW1fcmV3YXJkcw==',
        'YLtgEaeB64tT1h15BHRQrapGl34gI5MSlRJiwQV9EJA=',
        'DXk=',
        'PyR/',
        'ZSestA==',
        'U1VUSy0xNTk4OGE=',
    ],
};
