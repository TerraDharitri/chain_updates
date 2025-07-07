import { registerEnumType } from "@nestjs/graphql";

export enum MoaPairExchange {
  dharitrix = 'dharitrix',
  unknown = 'unknown'
}

registerEnumType(MoaPairExchange, {
  name: 'MoaPairExchange',
  description: 'MoaPairExchange object type.',
  valuesMap: {
    dharitrix: {
      description: 'dharitrix',
    },
    unknown: {
      description: 'unknown',
    },
  },
});
