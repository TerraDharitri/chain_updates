import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { AuthService } from '../services/auth.service';
import { RouteNamesEnum } from '../route-names.enum';
import { WidgetType } from '../../types/widget.types';
import { WidgetComponent } from '../components';
import {
  AccountComponent,
  PingPongRawComponent,
  SignMessageComponent,
} from '../widgets';

const WIDGETS: WidgetType[] = [
  {
    title: 'Account',
    widget: AccountComponent,
    description: 'Connected account details',
    reference: 'https://docs.multiversx.com/sdk-and-tools/sdk-dapp/#account',
  },
  {
    title: 'Ping & Pong (Manual)',
    widget: PingPongRawComponent,
    description:
      'Smart Contract interactions using manually formulated transactions',
    reference:
      'https://docs.multiversx.com/sdk-and-tools/indices/es-index-transactions/',
  },
  {
    title: 'Sign message',
    widget: SignMessageComponent,
    description: 'Message signing using the connected account',
    reference: 'https://docs.multiversx.com/sdk-and-tools/sdk-dapp/#account-1',
  },
];

@Component({
  selector: 'app-dashboard-page',
  standalone: true,
  imports: [CommonModule, WidgetComponent],
  templateUrl: './dashboard-page.component.html',
  styleUrls: ['./dashboard-page.component.css'],
})
export class DashboardPageComponent {
  widgets = WIDGETS;

  constructor(private router: Router, private authService: AuthService) {}

  trackByTitle(index: number, item: WidgetType): string {
    return item.title;
  }

  async logout() {
    await this.authService.logout();
    this.router.navigate([RouteNamesEnum.unlock]);
  }
}
