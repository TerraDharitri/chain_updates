import { Field, InputType, Int } from '@nestjs/graphql';
import { IsOptional, Matches } from 'class-validator';
import { drtConfig } from 'src/config';
import {
  REWA_OR_DCDT_TOKEN_RGX,
  DCDT_TOKEN_ERROR,
  NFT_IDENTIFIER_ERROR,
  NFT_IDENTIFIER_RGX,
  NUMERIC_ERROR,
  NUMERIC_RGX,
} from 'src/utils/constants';

@InputType()
export class BidActionArgs {
  @Field(() => Int)
  auctionId: number;

  @Matches(RegExp(NFT_IDENTIFIER_RGX), {
    message: NFT_IDENTIFIER_ERROR,
  })
  @Field(() => String)
  identifier: string;

  @Matches(RegExp(REWA_OR_DCDT_TOKEN_RGX), {
    message: DCDT_TOKEN_ERROR,
  })
  @Field(() => String)
  paymentTokenIdentifier: string = drtConfig.rewa;

  @Field(() => String)
  price: string;
}

@InputType()
export class BuySftActionArgs {
  @Field(() => Int)
  auctionId: number;

  @Matches(RegExp(NFT_IDENTIFIER_RGX), {
    message: NFT_IDENTIFIER_ERROR,
  })
  @Field(() => String)
  identifier: string;

  @Matches(RegExp(REWA_OR_DCDT_TOKEN_RGX), {
    message: DCDT_TOKEN_ERROR,
  })
  @Field(() => String)
  paymentTokenIdentifier: string = drtConfig.rewa;

  @Matches(RegExp(NUMERIC_RGX), { message: `Price ${NUMERIC_ERROR}` })
  @Field(() => String)
  price: string;

  @IsOptional()
  @Matches(RegExp(NUMERIC_RGX), { message: `Quantity ${NUMERIC_ERROR}` })
  @Field(() => String, { nullable: true })
  quantity: string;
}
