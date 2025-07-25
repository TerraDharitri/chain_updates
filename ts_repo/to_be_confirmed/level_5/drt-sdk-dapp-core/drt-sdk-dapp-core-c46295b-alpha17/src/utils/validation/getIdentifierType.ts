const dcdtParts = 2;
const nftParts = 3;

const defaultResult = {
  isDcdt: false,
  isNft: false,
  isRewa: false
};

export function getIdentifierType(identifier?: string): {
  isDcdt: boolean;
  isNft: boolean;
  isRewa: boolean;
} {
  const parts = identifier?.split('-').length;

  if (parts === dcdtParts) {
    return {
      ...defaultResult,
      isDcdt: true
    };
  }
  if (parts === nftParts) {
    return {
      ...defaultResult,
      isNft: true
    };
  }
  return {
    ...defaultResult,
    isRewa: true
  };
}
