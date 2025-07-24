import { Field, InputType } from '@nestjs/graphql';
import { drtConfig } from 'src/config';

@InputType()
export class TokenFilter {
  @Field(() => String)
  token: string = drtConfig.rewa;
  constructor(init?: Partial<TokenFilter>) {
    Object.assign(this, init);
  }
}
