import { VoteEvent } from '../governance.event';
import { rawVoteEvent } from '../mocks/mocked.raw.event';

describe('test vote event decoder', () => {
    it('should decode', () => {
        const voteEvent = new VoteEvent(
            rawVoteEvent,
        );
        expect(voteEvent.getTopics()).toEqual(
            {
                eventName: "upVoteCast",
                voter: "drt1p8lkt4t7gy8hwq90egjetyzkrfzcvruxy7dsqh4m3nefh94yys9smkrgt3",
                proposalId: 1,
                nrVotes: "9086692690078696219362",
                quorumUsed: "9086692690078696219362",
            });
    });
});
