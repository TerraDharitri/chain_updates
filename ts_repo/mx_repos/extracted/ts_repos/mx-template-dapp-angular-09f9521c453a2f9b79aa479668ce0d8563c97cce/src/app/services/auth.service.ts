import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { getAccountProvider } from '@multiversx/sdk-dapp/out/providers/helpers/accountProvider';
import { getIsLoggedIn } from '@multiversx/sdk-dapp/out/methods/account/getIsLoggedIn';
import { getStore } from '@multiversx/sdk-dapp/out/store/store';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private isLoggedInSubject = new BehaviorSubject<boolean>(false);
  public isLoggedIn$: Observable<boolean> = this.isLoggedInSubject.asObservable();
  private storeUnsubscribe?: () => void;

  constructor() {
    // Initialize login state - check actual login state from SDK
    this.checkLoginState();
    this.subscribeToStoreChanges();
  }

  private checkLoginState() {
    // Check the actual login state from SDK-DAPP
    const isLoggedIn = getIsLoggedIn();
    this.isLoggedInSubject.next(isLoggedIn);
  }

  private subscribeToStoreChanges() {
    // Subscribe to store changes to keep login state in sync
    const store = getStore();
    this.storeUnsubscribe = store.subscribe(() => {
      const isLoggedIn = getIsLoggedIn();
      if (this.isLoggedInSubject.value !== isLoggedIn) {
        this.isLoggedInSubject.next(isLoggedIn);
      }
    });
  }

  getIsLoggedIn(): boolean {
    return this.isLoggedInSubject.value;
  }

  setLoggedIn(loggedIn: boolean) {
    this.isLoggedInSubject.next(loggedIn);
  }

  async logout() {
    const provider = getAccountProvider();
    await provider.logout();
    this.setLoggedIn(false);
  }

  ngOnDestroy() {
    // Clean up subscription when service is destroyed
    if (this.storeUnsubscribe) {
      this.storeUnsubscribe();
    }
  }
} 