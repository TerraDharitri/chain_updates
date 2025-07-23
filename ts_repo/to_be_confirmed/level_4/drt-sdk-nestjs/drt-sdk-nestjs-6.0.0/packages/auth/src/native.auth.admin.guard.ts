import { Injectable, CanActivate, ExecutionContext, Inject } from '@nestjs/common';
import { ExecutionContextUtils, DrtnestConfigService, DRTNEST_CONFIG_SERVICE } from '@terradharitri/sdk-nestjs-common';

/**
 * This Guard allows only specific addresses to be authenticated.
 * 
 * The addresses are defined in the config file and are passed to the guard via the DrtnestConfigService.
 *
 * @return {boolean} `canActivate` returns true if the address is in the list of admins and uses a valid Native-Auth token.
 * 
 * @param {CachingService} CachingService - Dependency of `NativeAuthGuard`
 * @param {DrtnestConfigService} DrtnestConfigService - Dependency of `NativeAuthGuard`. Also used to get the list of admins (`getSecurityAdmins`). 
 */
@Injectable()
export class NativeAuthAdminGuard implements CanActivate {
  constructor(
    @Inject(DRTNEST_CONFIG_SERVICE)
    private readonly drtnestConfigService: DrtnestConfigService
  ) { }

  canActivate(context: ExecutionContext): boolean {
    const admins = this.drtnestConfigService.getSecurityAdmins();
    if (!admins) {
      return false;
    }

    const request = ExecutionContextUtils.getRequest(context);

    return admins.includes(request.nativeAuth.signerAddress);
  }
}
