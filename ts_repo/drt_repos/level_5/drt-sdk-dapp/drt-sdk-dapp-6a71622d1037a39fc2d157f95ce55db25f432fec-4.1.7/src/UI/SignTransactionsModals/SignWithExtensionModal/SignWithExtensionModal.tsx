import React from 'react';
import { SignModalPropsType } from 'types';

import {
  SignWaitingScreenModal,
  SignWaitingScreenModalPropsType
} from '../components';

export const SignWithExtensionModal = (props: SignModalPropsType) => {
  const description = props.error
    ? props.error
    : props.transactions?.length > 1
    ? 'Check your Dharitri Wallet Extension to sign the transactions'
    : 'Check your Dharitri Wallet Extension to sign the transaction';

  const waitingScreenProps: SignWaitingScreenModalPropsType = {
    ...props,
    description,
    title: 'Confirm on Dharitri DeFi Wallet'
  };

  return <SignWaitingScreenModal {...waitingScreenProps} />;
};
