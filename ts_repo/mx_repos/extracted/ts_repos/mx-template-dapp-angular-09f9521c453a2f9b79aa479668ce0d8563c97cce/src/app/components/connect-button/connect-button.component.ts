import { Component, Input } from '@angular/core';
import { Router } from '@angular/router';
import { UnlockPanelManager } from '@multiversx/sdk-dapp/out/managers/UnlockPanelManager';
import { RouteNamesEnum } from '../../route-names.enum';
import { ButtonComponent } from '../button/button.component';

@Component({
  selector: 'app-connect-button',
  standalone: true,
  imports: [ButtonComponent],
  templateUrl: './connect-button.component.html',
  styleUrls: ['./connect-button.component.css']
})
export class ConnectButtonComponent {
  @Input() className: string = 'inline-block rounded-lg px-3 py-2 text-center hover:no-underline my-0 bg-blue-600 text-white hover:bg-blue-700 mr-0';

  constructor(
    private router: Router
  ) {}

  login() {
    const unlockPanelManager = UnlockPanelManager.init({
      loginHandler: () => {
        this.router.navigate([RouteNamesEnum.dashboard]);
      },
      onClose: () => {}
    });
    unlockPanelManager.openUnlockPanel();
  }
} 