import { Injectable, CanActivate, ExecutionContext, Inject, Optional } from '@nestjs/common';
import { CacheService } from '@terradharitri/sdk-nestjs-cache';
import { MxnestConfigService, DRTNEST_CONFIG_SERVICE } from '@terradharitri/sdk-nestjs-common';
import { JwtAuthenticateGuard } from './jwt.authenticate.guard';
import { NativeAuthGuard } from './native.auth.guard';

@Injectable()
export class JwtOrNativeAuthGuard implements CanActivate {
  constructor(
    @Inject(DRTNEST_CONFIG_SERVICE) private readonly drtnestConfigService: MxnestConfigService,
    @Optional() private readonly cacheService?: CacheService,
  ) { }

  async canActivate(context: ExecutionContext): Promise<boolean> {
    const jwtGuard = new JwtAuthenticateGuard(this.drtnestConfigService);
    const nativeAuthGuard = new NativeAuthGuard(this.drtnestConfigService, this.cacheService);

    try {
      const result = await jwtGuard.canActivate(context);
      if (result) {
        return true;
      }
    } catch (error) {
      // do nothing
    }

    try {
      return await nativeAuthGuard.canActivate(context);
    } catch (error) {
      return false;
    }
  }
}
