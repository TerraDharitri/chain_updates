export const isRewaToken = (name?: string) => {
  if (!name) {
    return false;
  }
  return ['rewa', 'xrewa', 'wrewa'].includes(name?.toLowerCase());
};
