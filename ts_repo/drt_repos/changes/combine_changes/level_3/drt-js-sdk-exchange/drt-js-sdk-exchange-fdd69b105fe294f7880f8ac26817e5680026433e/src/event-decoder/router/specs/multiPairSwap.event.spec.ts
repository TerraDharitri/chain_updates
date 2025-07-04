import { rawMultiPairSwapEvent } from '../mocks/mocked.raw.event';
import { MultiPairSwapEvent } from '../multiPairSwap.event';

describe('test router events decoder', () => {
    it('should decode multi pair swap event', () => {
        const multiPairSwapEvent = new MultiPairSwapEvent(
            rawMultiPairSwapEvent,
        );
        expect(multiPairSwapEvent.toJSON()).toEqual({
            address:
                'drt1qqqqqqqqqqqqqpgqq66xk9gfr4esuhem3jru86wg5hvp33a62jpsh4nhal',
            identifier: 'multiPairSwap',
            caller: 'drt1qqqqqqqqqqqqqpgqsytkvnexypp7argk02l0rasnj57sxa542jps22gawh',
            block: 25574671,
            epoch: 1777,
            timestamp: 1749709740,
            tokenInID: 'WREWA-bd4d79',
            amountIn: '700000000000000000',
            tokenOutID: 'USDC-c76f1f',
            amountOut: '10567308',
            paymentsOut: [
                {
                    tokenIdentifier: 'USDC-c76f1f',
                    tokenNonce: 0,
                    amount: '10567308',
                },
            ],
        });
    });
});
