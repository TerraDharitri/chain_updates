import { dcdtAttributes } from './dcdt';

const baseUserTokenAttributes = `
  __typename
    valueUSD
`;

export const userDcdtAttributes = `
  ${dcdtAttributes}
  ${baseUserTokenAttributes}
`;
