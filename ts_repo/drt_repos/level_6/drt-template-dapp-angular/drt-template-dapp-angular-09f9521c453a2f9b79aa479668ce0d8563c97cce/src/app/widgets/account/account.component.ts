import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { OutputContainerComponent } from '../../components/output-container/output-container.component';
import { LabelComponent } from '../../components/label/label.component';
import { getAccount } from '@terradharitri/sdk-dapp/out/methods/account/getAccount';
import { getStore } from '@terradharitri/sdk-dapp/out/store/store';

@Component({
  selector: 'app-account',
  standalone: true,
  imports: [CommonModule, OutputContainerComponent, LabelComponent],
  templateUrl: './account.component.html',
  styleUrls: ['./account.component.css']
})
export class AccountComponent implements OnInit, OnDestroy {
  address: string = '';
  balance: string = '';
  shard: number = 0;
  nonce: number = 0;
  isLoading: boolean = true;
  storeUnsubscribe: (() => void) | undefined;

  async ngOnInit() {
    // Initial load
    this.updateAccount();

    // Subscribe to store changes
    const store = getStore();
    this.storeUnsubscribe = store.subscribe(() => {
      this.updateAccount();
    });
  }

  updateAccount() {
    const account = getAccount();
    
    if (account) {
      this.address = account.address;
      this.balance = account.balance.toString();
      this.shard = account.shard || 0;
      this.nonce = account.nonce || 0;
      this.isLoading = false;
    }
  }

  ngOnDestroy() {
    if (this.storeUnsubscribe) {
      this.storeUnsubscribe();
    }
  }
} 