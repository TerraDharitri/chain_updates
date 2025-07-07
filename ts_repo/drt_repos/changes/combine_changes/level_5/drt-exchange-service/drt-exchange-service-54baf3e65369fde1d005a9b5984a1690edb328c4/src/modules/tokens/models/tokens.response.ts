import { ObjectType } from '@nestjs/graphql';
import relayTypes from 'src/utils/relay.types';
import { DcdtToken } from './dcdtToken.model';

@ObjectType()
export class TokensResponse extends relayTypes<DcdtToken>(DcdtToken) {}
