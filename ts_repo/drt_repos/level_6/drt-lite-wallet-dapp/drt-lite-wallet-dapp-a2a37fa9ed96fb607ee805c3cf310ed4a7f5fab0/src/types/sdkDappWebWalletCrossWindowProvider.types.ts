import { WindowProviderResponseEnums as WindowResponseEnum } from '@terradharitri/sdk-web-wallet-cross-window-provider/out/enums';

export {
  WindowProviderRequestEnums,
  WindowProviderResponseEnums,
  SignMessageStatusEnum
} from '@terradharitri/sdk-web-wallet-cross-window-provider/out/enums';
import { ReplyWithPostMessageObjectType as OriginalReplyWithPostMessageObjectType } from '@terradharitri/sdk-web-wallet-cross-window-provider/out/types';
export type { RequestMessageType } from '@terradharitri/sdk-web-wallet-cross-window-provider/out/types';
export type { ReplyWithRedirectType } from '@terradharitri/sdk-js-web-wallet-io/out/replyToDapp/replyWithRedirect';

type ExtendedLoginResponseType = Omit<
  OriginalReplyWithPostMessageObjectType[WindowResponseEnum.loginResponse],
  'signature'
> & {
  name?: string; // allow the presence of name field for extension
  signature?: string; // make signature optional
}; // ⬇️ propagate the above changes down the chain

type ExtendedReplyWithPostMessageObjectType = Omit<
  OriginalReplyWithPostMessageObjectType,
  WindowResponseEnum.loginResponse
> & {
  [WindowResponseEnum.loginResponse]: ExtendedLoginResponseType;
};

type ReplyWithPostMessagePayloadType<
  K extends keyof ExtendedReplyWithPostMessageObjectType
> = {
  data?: ExtendedReplyWithPostMessageObjectType[K];
  error?: string;
};

// this is overwriting the original
export type ExtendedReplyWithPostMessageType = {
  [K in keyof ExtendedReplyWithPostMessageObjectType]: {
    type: K;
    payload: ReplyWithPostMessagePayloadType<K>;
  };
}[keyof ExtendedReplyWithPostMessageObjectType];
