import { IWrapAbiService } from '../services/interfaces';
import { WrapAbiService } from '../services/wrap.abi.service';

export class WrapAbiServiceMock implements IWrapAbiService {
    async wrappedRewaTokenID(): Promise<string> {
        return 'WREWA-123456';
    }
}

export const WrapAbiServiceProvider = {
    provide: WrapAbiService,
    useClass: WrapAbiServiceMock,
};
