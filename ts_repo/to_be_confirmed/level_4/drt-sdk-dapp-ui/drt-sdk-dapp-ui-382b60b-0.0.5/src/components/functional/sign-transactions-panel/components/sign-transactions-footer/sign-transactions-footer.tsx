import { Component, Fragment, h, State } from '@stencil/core';
import { DataTestIdsEnum } from 'constants/dataTestIds.enum';

import state from '../../signTransactionsPanelStore';

const signTransactionsFooterClasses: Record<string, string> = {
  buttonTooltip: 'drt:absolute drt:top-0 drt:h-12 drt:left-0 drt:right-0',
};

@Component({
  tag: 'drt-sign-transactions-footer',
  styleUrl: 'sign-transactions-footer.scss',
  shadow: true,
})
export class SignTransactionsFooter {
  @State() awaitsExternalConfirmation: boolean = false;

  render() {
    const { onCancel, onBack, onNext, onConfirm } = state;
    const { currentIndex, currentIndexToSign, needsSigning, username, address, explorerLink } = state.commonData;

    const isFirstTransaction = currentIndex === 0;
    const currentIndexNeedsSigning = currentIndex === currentIndexToSign;
    const currentIndexCannotBeSignedYet = currentIndex > currentIndexToSign;
    const showForwardAction = currentIndexNeedsSigning || currentIndexCannotBeSignedYet;

    return (
      <div class="sign-transactions-footer">
        <div class="sign-transactions-footer-buttons">
          <div class="sign-transactions-footer-button-wrapper cancel">
            <button
              data-testid={isFirstTransaction ? DataTestIdsEnum.signCancelBtn : DataTestIdsEnum.signBackBtn}
              onClick={isFirstTransaction ? onCancel : onBack}
              class={{
                'sign-transactions-footer-button': true,
                'cancel': !currentIndexCannotBeSignedYet,
                'highlighted': currentIndexCannotBeSignedYet,
              }}
            >
              {isFirstTransaction ? 'Cancel' : 'Back'}
            </button>
          </div>

          <div class="sign-transactions-footer-button-wrapper confirm">
            {currentIndexCannotBeSignedYet && (
              <div
                class="sign-transactions-footer-button-tooltip-wrapper"
                onClick={(event: MouseEvent) => event.stopPropagation()}
              >
                <drt-tooltip
                  trigger={
                    <div
                      class={{
                        'sign-transactions-footer-button-tooltip': true,
                        [signTransactionsFooterClasses.buttonTooltip]: true,
                      }}
                    />
                  }
                >
                  {needsSigning ? (
                    <Fragment>
                      You cannot sign this transaction yet, <br /> please go back and sign consecutively.
                    </Fragment>
                  ) : (
                    <Fragment>
                      You cannot confirm this transaction yet, <br />
                      please go back and confirm consecutively.
                    </Fragment>
                  )}
                </drt-tooltip>
              </div>
            )}

            <button
              data-testid={DataTestIdsEnum.signNextTransactionBtn}
              onClick={showForwardAction ? onConfirm : onNext}
              class={{
                'sign-transactions-footer-button': true,
                'highlighted': true,
                'disabled': currentIndexCannotBeSignedYet,
              }}
            >
              {showForwardAction ? (
                <span class="sign-transactions-footer-button-label">{needsSigning ? 'Sign' : 'Confirm'}</span>
              ) : (
                <span class="sign-transactions-footer-button-label">Next</span>
              )}

              {showForwardAction ? (
                <span
                  class={{ 'sign-transactions-footer-button-icon': true, 'lighter': currentIndexCannotBeSignedYet }}
                >
                  {needsSigning ? <drt-pencil-icon /> : <drt-check-icon />}
                </span>
              ) : (
                <span class="sign-transactions-footer-button-icon">
                  <drt-arrow-right-icon />
                </span>
              )}
            </button>
          </div>
        </div>

        <div class="sign-transactions-footer-identity">
          <div class="sign-transactions-footer-identity-label">Sign in with</div>

          {username && (
            <div class="sign-transactions-footer-identity-username">
              <span class="sign-transactions-footer-identity-username-prefix">@</span>
              <span class="sign-transactions-footer-identity-username-text">{username}</span>
            </div>
          )}

          {!username && address && <drt-trim-text text={address} class="sign-transactions-footer-identity-address" />}

          <drt-copy-button
            text={username ?? address}
            class="sign-transactions-footer-identity-copy"
            iconClass="sign-transactions-footer-identity-copy-icon"
          />
          <drt-explorer-link
            link={explorerLink}
            class="sign-transactions-footer-identity-explorer"
            iconClass="sign-transactions-footer-identity-explorer-icon"
          />
        </div>
      </div>
    );
  }
}
