import { registerEnumType } from "@nestjs/graphql";

export enum DcdtDataSource {
  gateway = 'gateway',
  elastic = 'elastic'
}

registerEnumType(DcdtDataSource, {
  name: 'DcdtDataSource',
  description: 'DCDT data source.',
  valuesMap: {
    gateway: {
      description: 'Gateway data source.',
    },
    elastic: {
      description: 'Elastic data source.',
    },
  },
});
