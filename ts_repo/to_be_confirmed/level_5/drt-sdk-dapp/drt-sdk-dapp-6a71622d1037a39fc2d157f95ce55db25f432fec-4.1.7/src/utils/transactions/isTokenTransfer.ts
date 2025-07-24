export function isTokenTransfer({
  tokenId,
  drtLabel
}: {
  tokenId: string | undefined;
  drtLabel: string;
}) {
  return Boolean(tokenId && tokenId !== drtLabel);
}
