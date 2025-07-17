export const trimUsernameDomain = (username?: string) => {
  if (!username) {
    return;
  }

  const numbatSuffixExists = username.lastIndexOf('.numbat') > 0;
  const trimmedPartBeforeLastDot = numbatSuffixExists
    ? username.substring(0, username.lastIndexOf('.'))
    : username;

  return trimmedPartBeforeLastDot;
};
