import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import {
  faPen,
  faBroom,
  faRotateRight,
} from '@fortawesome/free-solid-svg-icons';
import { OutputContainerComponent } from '../../components/output-container/output-container.component';
import { LabelComponent } from '../../components/label/label.component';
import { ButtonComponent } from '../../components/button/button.component';
import { getAccount } from '@terradharitri/sdk-dapp/out/methods/account/getAccount';
import { getAccountProvider } from '@terradharitri/sdk-dapp/out/providers/helpers/accountProvider';
import { Message } from '@terradharitri/sdk-core/out';

@Component({
  selector: 'app-sign-message',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    FontAwesomeModule,
    OutputContainerComponent,
    LabelComponent,
    ButtonComponent,
  ],
  templateUrl: './sign-message.component.html',
  styleUrls: ['./sign-message.component.css'],
})
export class SignMessageComponent {
  message: string = '';
  signedMessage: any = null;
  state: 'pending' | 'success' | 'error' = 'pending';
  signature: string = '';
  address: string = '';

  // FontAwesome icons
  faPen = faPen;
  faBroom = faBroom;
  faRotateRight = faRotateRight;

  async handleSubmit() {
    try {
      // Get the account data from the store
      const account = getAccount();
      this.address = account.address || '';

      if (!this.address) {
        console.error('No account address found');
        this.state = 'error';
        return;
      }

      console.log('Signing message:', this.message);

      // Create a Message object from the message
      const messageToSign = new Message({
        data: Buffer.from(this.message, 'utf8'),
      });

      // Get the account provider and sign the message
      const provider = getAccountProvider();
      const signedMessage = await provider.signMessage(messageToSign);

      if (!signedMessage) {
        console.error('No signed message found');
        this.state = 'error';
        return;
      }

      // Extract the signature from the signed message
      this.signature = signedMessage.signature
        ? Buffer.from(signedMessage.signature).toString('hex')
        : '';
      this.state = 'success';
      this.signedMessage = signedMessage;
      this.message = '';
    } catch (error) {
      console.error('Error signing message:', error);
      this.state = 'error';
    }
  }

  handleClear(event: MouseEvent) {
    event.preventDefault();
    event.stopPropagation();
    this.signature = '';
    this.signedMessage = null;
    this.state = 'pending';
    this.message = '';
  }

  get encodedMessage(): string {
    if (!this.signedMessage) return '';
    return (
      '0x' +
      Array.from(this.signedMessage.data, (byte: number) =>
        byte.toString(16).padStart(2, '0')
      ).join('')
    );
  }

  get decodedMessage(): string {
    if (!this.signedMessage) return '';
    return new TextDecoder().decode(this.signedMessage.data);
  }
}
