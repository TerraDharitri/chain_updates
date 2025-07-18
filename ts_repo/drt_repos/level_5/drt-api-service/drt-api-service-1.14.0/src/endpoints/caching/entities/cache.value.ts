import { Constants } from "@terradharitri/sdk-nestjs-common";

export class CacheValue {
  value?: string;
  ttl: number = Constants.oneSecond() * 6;
}
