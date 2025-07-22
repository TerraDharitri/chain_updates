import { DcdtTokenPayment } from '../dcdt.token.payment';
import {
    DcdtTokenPaymentAttributesNew,
    DcdtTokenPaymentAttributesOld,
} from '../mocks/dcdt.token.payment.mock';

describe('test dcdt token payment attributes decoder', () => {
    it('should decode old dcdt token payment', () => {
        const dcdtTokenPayment = DcdtTokenPayment.fromAttributes(
            DcdtTokenPaymentAttributesOld,
        );

        expect(dcdtTokenPayment).toEqual(
            new DcdtTokenPayment({
                tokenIdentifier: 'MOA-dc289c',
                tokenNonce: 0,
                amount: '365359339001228979577216',
            }),
        );
    });

    it('should decode dcdt token payment', () => {
        const dcdtTokenPayment = DcdtTokenPayment.fromAttributes(
            DcdtTokenPaymentAttributesNew,
        );

        expect(dcdtTokenPayment).toEqual(
            new DcdtTokenPayment({
                tokenIdentifier: 'MOA-dc289c',
                tokenNonce: 0,
                amount: '365359339001228979577216',
            }),
        );
    });
});
