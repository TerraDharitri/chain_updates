import { Module } from '@nestjs/common';
import { OwnersService } from './owners.service';
import { DrtCommunicationModule } from 'src/common';
import { OwnersResolver } from './owners.resolver';
import { AccountsProvider } from '../account-stats/loaders/accounts.loader';
import { AccountsRedisHandler } from '../account-stats/loaders/accounts.redis-handler';

@Module({
  providers: [OwnersService, OwnersResolver, AccountsRedisHandler, AccountsProvider],
  imports: [DrtCommunicationModule],
  exports: [OwnersService, AccountsRedisHandler, AccountsProvider],
})
export class OwnersModuleGraph {}
