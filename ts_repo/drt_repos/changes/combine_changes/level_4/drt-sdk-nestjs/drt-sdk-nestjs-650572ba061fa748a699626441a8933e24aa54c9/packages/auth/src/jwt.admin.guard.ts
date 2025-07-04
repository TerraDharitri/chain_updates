import { Injectable, CanActivate, ExecutionContext, Inject } from '@nestjs/common';
import { ExecutionContextUtils, MxnestConfigService, DRTNEST_CONFIG_SERVICE } from '@terradharitri/sdk-nestjs-common';

@Injectable()
export class JwtAdminGuard implements CanActivate {
  constructor(
    @Inject(DRTNEST_CONFIG_SERVICE)
    private readonly drtnestConfigService: MxnestConfigService
  ) { }

  // eslint-disable-next-line require-await
  async canActivate(
    context: ExecutionContext,
  ): Promise<boolean> {


    const admins = this.drtnestConfigService.getSecurityAdmins();
    if (!admins) {
      return false;
    }

    const request = ExecutionContextUtils.getRequest(context);

    return admins.includes(request.jwt.address);
  }
}
