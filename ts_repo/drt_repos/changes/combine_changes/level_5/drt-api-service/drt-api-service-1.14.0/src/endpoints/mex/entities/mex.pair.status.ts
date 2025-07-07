import { registerEnumType } from "@nestjs/graphql";

export enum MoaPairStatus {
  active = 'Active',
  inactive = 'Inactive',
  paused = 'Paused',
  partial = 'Partial',
}

registerEnumType(MoaPairStatus, {
  name: 'MoaPairStatus',
  description: 'MoaPairStatus object type.',
  valuesMap: {
    active: {
      description: 'Active state.',
    },
    inactive: {
      description: 'Inactive state.',
    },
    paused: {
      description: 'Pause state.',
    },
    partial: {
      description: 'Partial state.',
    },
  },
});
