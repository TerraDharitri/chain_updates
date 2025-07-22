import { confirmationDialogTag } from './constants';
import { getStyles } from './getStyles';
import { PopupConsentModel } from './PopupConsent.model';

export class PopupConsent extends HTMLElement implements PopupConsentModel {
  public walletUrl = '';
  public identifier: string = confirmationDialogTag;
  public onCancel = () => {
    console.log('onCancel');
  };
  public onConfirm = () => {
    console.log('onConfirm');
  };

  constructor() {
    super();
  }

  attributeChangedCallback(name: string, oldValue: string, newValue: string) {
    if (oldValue !== newValue) {
      (this as any)[name.replace(/-./g, (x) => x[1].toUpperCase())] = newValue;
      this.render();
    }
  }

  public events = {
    confirm: `${this.identifier}-confirm`,
    cancel: `${this.identifier}-cancel`,
    mounted: `${this.identifier}-mounted`
  };

  private handleConfirmEvent({ type }: CustomEvent) {
    switch (type) {
      case this.events.cancel:
        return this.onCancel();
      case this.events.confirm:
        return this.onConfirm();
      default:
        break;
    }
  }

  private dispatchAction = (type: string) => {
    const event = new CustomEvent(type, {
      bubbles: true, // Allow the event to bubble up through the DOM
      composed: true // Allow the event to cross the shadow DOM boundary
    });
    this.dispatchEvent(event);
  };

  render() {
    this.innerHTML = `
        <style>
          ${getStyles(confirmationDialogTag)}
        </style>
        <div id="${confirmationDialogTag}">
            <div class="content">
                <div class="body">
                    <div class="title">Confirm on DharitrI Wallet</div>
                    <div class="subtitle">Continue to ${this.walletUrl}</div>
                    <div class="actions-container">
                        <button
                            class="button"
                            data-testid="${this.events.cancel}-btn"
                            id="${this.events.cancel}-btn"
                        >
                            Cancel
                        </button>
                        <button
                            class="button btn-proceed"
                            data-testid="${this.events.confirm}-btn"
                            id="${this.events.confirm}-btn"
                        >
                            Continue →
                        </button>
                    </div>
                </div>
            </div>
        </div>
      `;
  }

  toggleEvents(action: 'removeEventListener' | 'addEventListener') {
    [this.events.cancel, this.events.confirm].forEach((event) => {
      this.querySelector(`#${event}-btn`)?.[action]('click', () =>
        this.dispatchAction(event)
      );
      this[action](event, this.handleConfirmEvent as EventListener);
    });
  }

  connectedCallback() {
    this.render();
    this.toggleEvents('addEventListener');
  }

  disconnectedCallback() {
    this.toggleEvents('removeEventListener');
  }
}
