import { nativeAuth } from '@terradharitri/sdk-dapp/services/nativeAuth/nativeAuth';
import { replyToDapp as originalReplyToDapp } from '@terradharitri/sdk-js-web-wallet-io/out/replyToDapp/replyToDapp';
import { ExtendedReplyWithPostMessageType, ReplyWithRedirectType } from 'types';

export { getRewaLabel } from '@terradharitri/sdk-dapp/utils/network/getRewaLabel';
export { getTransactions } from '@terradharitri/sdk-dapp/apiCalls/transactions/getTransactions';
export { sendTransactions } from '@terradharitri/sdk-dapp/services/transactions/sendTransactions';
export { sendBatchTransactions as sendBatchTransactionsSdkDapp } from '@terradharitri/sdk-dapp/services/transactions/sendBatchTransactions';
export { refreshAccount } from '@terradharitri/sdk-dapp/utils/account/refreshAccount';
export { logout } from '@terradharitri/sdk-dapp/utils/logout';
export { signTransactions } from '@terradharitri/sdk-dapp/services/transactions/signTransactions';
export { trimUsernameDomain } from '@terradharitri/sdk-dapp/hooks/account/helpers';
export { getAccount } from '@terradharitri/sdk-dapp/utils/account/getAccount';
export { getAddress } from '@terradharitri/sdk-dapp/utils/account/getAddress';
export { newTransaction } from '@terradharitri/sdk-dapp/models';
export { useLoginService } from '@terradharitri/sdk-dapp/hooks/login/useLoginService';
export { decodeNativeAuthToken } from '@terradharitri/sdk-dapp/services/nativeAuth/helpers/decodeNativeAuthToken';
export { getIsNativeAuthSingingForbidden } from '@terradharitri/sdk-dapp/services/nativeAuth/helpers/getIsNativeAuthSingingForbidden';
export { decodeLoginToken } from '@terradharitri/sdk-dapp/services/nativeAuth/helpers/decodeLoginToken';
export { getWebviewToken } from '@terradharitri/sdk-dapp/utils/account/getWebviewToken';
export { getAccountProviderType } from '@terradharitri/sdk-dapp/utils/account/getAccountProviderType';

const { getToken } = nativeAuth();
export { getToken };
export { loginWithExternalProvider } from '@terradharitri/sdk-dapp/utils/account/loginWithExternalProvider';
export { addressIsValid } from '@terradharitri/sdk-dapp/utils/account/addressIsValid';
export { getInterpretedTransaction } from '@terradharitri/sdk-dapp/utils/transactions/getInterpretedTransaction';
export { formatAmount } from '@terradharitri/sdk-dapp/utils/operations/formatAmount';
export { parseAmount } from '@terradharitri/sdk-dapp/utils/operations/parseAmount';
export { getIsProviderEqualTo } from '@terradharitri/sdk-dapp/utils/account/getIsProviderEqualTo';
export { removeTransactionsToSign } from '@terradharitri/sdk-dapp/services/transactions/clearTransactions';
export { extractSessionId } from '@terradharitri/sdk-dapp/hooks/transactions/helpers/extractSessionId';
export { checkIsValidSender } from '@terradharitri/sdk-dapp/hooks/transactions/helpers/checkIsValidSender';
export { useAddressScreens } from '@terradharitri/sdk-dapp/hooks/login/useAddressScreens';
export { useBatchTransactionsTracker } from '@terradharitri/sdk-dapp/hooks/transactions/batch/tracker/useBatchTransactionsTracker';
export { useCheckBatch } from '@terradharitri/sdk-dapp/hooks/transactions/batch/tracker/useCheckBatch';
export { useGetAccount } from '@terradharitri/sdk-dapp/hooks/account/useGetAccount';
export { useGetAccountInfo } from '@terradharitri/sdk-dapp/hooks/account/useGetAccountInfo';
export { useGetAccountProvider } from '@terradharitri/sdk-dapp/hooks/account/useGetAccountProvider';
export { useGetActiveTransactionsStatus } from '@terradharitri/sdk-dapp/hooks/transactions/useGetActiveTransactionsStatus';
export { useGetIsLoggedIn } from '@terradharitri/sdk-dapp/hooks/account/useGetIsLoggedIn';
export { useGetIsWalletConnectV2Initialized } from '@terradharitri/sdk-dapp/hooks/account/useGetIsWalletConnectV2Initialized';
export { useGetLastPendingTransactionHash } from '@terradharitri/sdk-dapp/hooks/transactions/useGetLastPendingTransactionHash';
export { useGetLastSignedMessageSession } from '@terradharitri/sdk-dapp/hooks/signMessage/useGetLastSignedMessageSession';
export { useGetSignMessageSession } from '@terradharitri/sdk-dapp/hooks/signMessage/useGetSignMessageSession';
export { useGetLoginInfo } from '@terradharitri/sdk-dapp/hooks/account/useGetLoginInfo';
export { useGetNetworkConfig } from '@terradharitri/sdk-dapp/hooks/useGetNetworkConfig';
export { useGetPendingTransactions } from '@terradharitri/sdk-dapp/hooks/transactions/useGetPendingTransactions';
export { useGetSignMessageInfoStatus } from '@terradharitri/sdk-dapp/hooks/signMessage/useGetSignedMessageStatus';
export { useGetSignTransactionsError } from '@terradharitri/sdk-dapp/hooks/transactions/useGetSignTransactionsError';
export { useGetSignedTransactions } from '@terradharitri/sdk-dapp/hooks/transactions/useGetSignedTransactions';
export { useSendBatchTransactions } from '@terradharitri/sdk-dapp/hooks/transactions/batch/useSendBatchTransactions';
export { useSignMessage } from '@terradharitri/sdk-dapp/hooks/signMessage/useSignMessage';
export { useSignTransactions } from '@terradharitri/sdk-dapp/hooks/transactions/useSignTransactions';
export { useSignTransactionsCommonData } from '@terradharitri/sdk-dapp/hooks/transactions/useSignTransactionsCommonData';
export { useSignTransactionsWithDevice } from '@terradharitri/sdk-dapp/hooks/transactions/useSignTransactionsWithDevice';
export { useSignTransactionsWithLedger } from '@terradharitri/sdk-dapp/hooks/transactions/useSignTransactionsWithLedger';
export { useTrackTransactionStatus } from '@terradharitri/sdk-dapp/hooks/transactions/useTrackTransactionStatus';
export { verifyMessage } from '@terradharitri/sdk-dapp/hooks/signMessage/verifyMessage';
export { useTransactionsTracker } from '@terradharitri/sdk-dapp/hooks/transactions/useTransactionsTracker';
export { useGetAccountFromApi } from '@terradharitri/sdk-dapp/apiCalls/accounts/useGetAccountFromApi';
export {
  deleteTransactionToast,
  removeAllSignedTransactions,
  removeAllTransactionsToSign
} from '@terradharitri/sdk-dapp/services/transactions/clearTransactions';
export {
  setTransactionsDisplayInfoState,
  setTransactionsToSignedState
} from '@terradharitri/sdk-dapp/services/transactions/updateSignedTransactions';
export { sendBatchTransactions } from '@terradharitri/sdk-dapp/services/transactions/sendBatchTransactions';
export { useAxiosInterceptorContext } from '@terradharitri/sdk-dapp/wrappers/AxiosInterceptorContext';
export { storage } from '@terradharitri/sdk-dapp/utils/storage';
export { addNewCustomToast } from '@terradharitri/sdk-dapp/utils/toasts';
export {
  maxDecimals,
  stringIsFloat,
  stringIsInteger
} from '@terradharitri/sdk-dapp/utils/validation';

export const replyToDapp: (
  props: {
    callbackUrl: string;
    webwiewApp?: HTMLIFrameElement | null;
    postMessageData?: ExtendedReplyWithPostMessageType;
    transactionData?: ReplyWithRedirectType['transactionData'];
  },
  extensionReplyToDapp?: (props: ExtendedReplyWithPostMessageType) => void
) => void = originalReplyToDapp as any; // use as any to allow extending input params
