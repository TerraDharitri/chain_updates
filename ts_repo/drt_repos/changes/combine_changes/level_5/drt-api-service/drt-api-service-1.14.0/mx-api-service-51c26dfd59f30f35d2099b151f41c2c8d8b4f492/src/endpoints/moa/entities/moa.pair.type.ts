import { registerEnumType } from "@nestjs/graphql";

export enum MoaPairType {
  core = 'core',
  community = 'community',
  ecosystem = 'ecosystem',
  experimental = 'experimental',
  unlisted = 'unlisted',
}

registerEnumType(MoaPairType, {
  name: 'MoaPairType',
  description: 'MoaPairType object type.',
  valuesMap: {
    core: {
      description: 'Core Type.',
    },
    community: {
      description: 'Community Type.',
    },
    ecosystem: {
      description: 'Ecosystem Type.',
    },
    experimental: {
      description: 'Experimental Type.',
    },
    unlisted: {
      description: 'Unlisted Type.',
    },
  },
});
