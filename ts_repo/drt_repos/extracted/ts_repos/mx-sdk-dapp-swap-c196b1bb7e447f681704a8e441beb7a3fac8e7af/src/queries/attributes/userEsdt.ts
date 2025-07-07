import { esdtAttributes } from './esdt';

const baseUserTokenAttributes = `
  __typename
    valueUSD
`;

export const userEsdtAttributes = `
  ${esdtAttributes}
  ${baseUserTokenAttributes}
`;
