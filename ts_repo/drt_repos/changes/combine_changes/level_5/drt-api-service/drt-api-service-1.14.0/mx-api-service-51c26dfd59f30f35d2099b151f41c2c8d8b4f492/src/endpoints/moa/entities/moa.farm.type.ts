import { registerEnumType } from "@nestjs/graphql";

export enum MoaFarmType {
  standard = 'standard',
  metastaking = 'metastaking',
}

registerEnumType(MoaFarmType, {
  name: 'MoaFarmType',
  description: 'MoaFarmType object type.',
  valuesMap: {
    standard: {
      description: 'Standard type.',
    },
    metastaking: {
      description: 'Metastaking type.',
    },
  },
});
