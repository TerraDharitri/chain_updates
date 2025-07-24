export interface DrtnestConfigService {
  getSecurityAdmins(): string[];

  getJwtSecret(): string;

  getApiUrl(): string;

  getNativeAuthMaxExpirySeconds(): number;

  getNativeAuthAcceptedOrigins(): string[];
}
