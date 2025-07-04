import { Component } from '@angular/core';
import { ButtonComponent } from '../button/button.component';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { faBell } from '@fortawesome/free-solid-svg-icons';
import { NotificationsFeedManager } from '@multiversx/sdk-dapp/out/managers/NotificationsFeedManager';

@Component({
  selector: 'app-notifications-button',
  standalone: true,
  imports: [ButtonComponent, FontAwesomeModule],
  templateUrl: './notifications-button.component.html',
  styleUrls: ['./notifications-button.component.css']
})
export class NotificationsButtonComponent {
  faBell = faBell;

  async handleOpenNotificationsFeed() {
    await NotificationsFeedManager.getInstance().openNotificationsFeed();
  }
} 