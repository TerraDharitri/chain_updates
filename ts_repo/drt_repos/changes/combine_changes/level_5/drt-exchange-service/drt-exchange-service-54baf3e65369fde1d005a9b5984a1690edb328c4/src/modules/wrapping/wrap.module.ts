import { Module } from '@nestjs/common';
import { MXCommunicationModule } from '../../services/dharitri-communication/drt.communication.module';
import { TokenModule } from '../tokens/token.module';
import { WrapAbiService } from './services/wrap.abi.service';
import { WrapTransactionsService } from './services/wrap.transactions.service';
import { WrapResolver } from './wrap.resolver';
import { WrapService } from './services/wrap.service';

@Module({
    imports: [MXCommunicationModule, TokenModule],
    providers: [
        WrapService,
        WrapAbiService,
        WrapTransactionsService,
        WrapResolver,
    ],
    exports: [WrapAbiService, WrapService, WrapTransactionsService],
})
export class WrappingModule {}
