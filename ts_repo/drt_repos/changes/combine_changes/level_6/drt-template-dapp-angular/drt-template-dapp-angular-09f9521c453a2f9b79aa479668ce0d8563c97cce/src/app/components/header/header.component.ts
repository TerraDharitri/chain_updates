import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../services/auth.service';
import { EnvironmentService } from '../../services/environment.service';
import { RouteNamesEnum } from '../../route-names.enum';
import { ButtonComponent } from '../button/button.component';
import { MxLinkComponent } from '../drt-link/drt-link.component';
import { NotificationsButtonComponent } from '../notifications-button/notifications-button.component';
import { ConnectButtonComponent } from '../connect-button/connect-button.component';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [CommonModule, ButtonComponent, MxLinkComponent, NotificationsButtonComponent, ConnectButtonComponent],
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css']
})
export class HeaderComponent {
  routeNames = RouteNamesEnum;

  constructor(
    private authService: AuthService,
    private environmentService: EnvironmentService,
    private router: Router
  ) {}

  get isLoggedIn(): boolean {
    return this.authService.getIsLoggedIn();
  }

  get environment(): string {
    return this.environmentService.environment;
  }

  async handleLogout() {
    await this.authService.logout();
    this.router.navigate([RouteNamesEnum.home]);
  }
} 