import { Injectable } from '@angular/core';
import { CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { RouteNamesEnum } from '../route-names.enum';

@Injectable({
  providedIn: 'root'
})
export class AuthRedirectGuard implements CanActivate {
  constructor(
    private authService: AuthService,
    private router: Router
  ) {}

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): boolean {
    const isLoggedIn = this.authService.getIsLoggedIn();
    const requireAuth = route.data?.['requireAuth'] || false;

    if (isLoggedIn && !requireAuth) {
      this.router.navigate([RouteNamesEnum.dashboard]);
      return false;
    }

    if (!isLoggedIn && requireAuth) {
      this.router.navigate([RouteNamesEnum.home]);
      return false;
    }

    return true;
  }
} 