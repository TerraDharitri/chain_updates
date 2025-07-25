export function isTokenTransfer({
  tokenId,
  rewaLabel
}: {
  tokenId: string | undefined;
  rewaLabel: string;
}) {
  return Boolean(tokenId && tokenId !== rewaLabel);
}
