import { registerEnumType } from "@nestjs/graphql";

export enum MoaPairState {
  active = 'active',
  inactive = 'inactive',
  paused = 'paused',
  partial = 'partial',
}

registerEnumType(MoaPairState, {
  name: 'MoaPairState',
  description: 'MoaPairState object type.',
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
