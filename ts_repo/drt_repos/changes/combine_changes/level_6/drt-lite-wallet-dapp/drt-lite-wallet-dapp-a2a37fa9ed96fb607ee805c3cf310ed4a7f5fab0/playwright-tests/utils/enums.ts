/* eslint-disable quotes */
export enum AssertionEnum {
  contain = 'contain',
  include = 'include',
  beChecked = 'be.checked',
  beDisabled = 'be.disabled',
  haveValue = 'have.value'
}

export enum ConfirmMsgEnum {
  pending = 'Processing transactions',
  confirmed = 'Transactions successful'
}

export enum ErrorMessagesEnum {
  invalidHerotag = 'Invalid herotag',
  required = 'Required',
  invalidNumber = 'Invalid number',
  insufficientFunds = 'Insufficient funds',
  onlyDigitsAllowed = 'Only digits and one . allowed',
  tokenRequired = 'Token required',
  minAmountErr = 'Minimum amount: 0.0005',
  wrongCode = 'wrong code'
}

export enum AccountStatesEnum {
  // Use different wallets for the same state to bypass the cooldown period for the QR generator.

  unGuardedAccount = 'drt17rp8l3waauynrhtw0233qf3kcwxxv9e8c0yhkk5l72jrndjm9j4qykfhpc',
  unGuardedAccount2 = 'drt1r97szd9zxh526zptnryl690wc93rfnm5asgs7lznf6laywfsvajqapsppj',
  unGuardedAccount3 = 'drt13fntl99g2ft78cljfflj6cmhsyz90lf70fek4as2prrspnl3ca8sk83ljk',
  unGuardedAccount4 = 'drt1vph98qdtfne4wnax0r3saqee94x2ughzrpnvhgaz32qm4vdhnm8qad7c9a',
  unGuardedAccount6 = 'drt1ckgy909lrfmyrk9k9efq8dp3mxnzx72wukmpahryx0kyfc85l3csrht9fr',
  unGuardedAccount7 = 'drt1c5alcvu8smr825ac3qp9es5m8x3fewhm9rf8z0jenvu7dr8szk3shgrkeg',
  unGuardAccount8 = 'drt1rqjjr44l25fycfj6v8nazpjwmd6jy6ffsettrn5re3qhmt26zt3s6awle7',
  unGuardedAccount8 = 'drt16xyzh38893pgldq8m2v6xyqf79xn7yq3wdv7e3ang6mdrs8dm7zqygspx8',
  unGuardAccount9 = 'drt1vrygm0hl9dtakvmvnev90zer2xlckqfed42ar2y6vv4camrkvy3qwnmhd8',
  readyToGuardAccount = 'drt1693xsqq6d7dqgd4x7wnr2a5tg4rfayw0nh393nu5umhqe7e229fqakpgq4',
  readyToGuardAccount1 = 'drt1enj49v5xvksfv3asc2lzajg9tykxm5e9fs62vgfs4wafd7nssrfs2e80f9',
  guardedAccount = 'drt1knjjk9nvl6tey3k0w0wfq964zpfu3gdm84ex6nqmz7plvq69je5s422uqh'
}

export enum RoutesEnum {
  send = '/send',
  guardian = '/guardian',
  guardianRegister = '/guardian-register',
  dashboard = '/dashboard',
  staking = '/stake',
  chooseProvider = '/stake/choose-provider',
  chooseAmount = '/stake/choose-amount',
  claimRewards = '/stake/claim-rewards',
  sign = '/sign',
  signMessage = '/sign-message',
  nft = '/nft',
  hook2fa = '/hook/2fa',
  unlock = '/unlock'
}

export enum KeystoreIndexEnum {
  keystore1 = '1',
  keystore2 = '2'
}

export enum GlobalSelectorsEnum {
  sendTrxBtn = 'sendTrxBtn',
  modalTitle = 'modalTitle',
  continueBtn = 'continueBtn',
  amount = 'amount',
  confirmData = 'confirmData',
  transactionStatusSubtitle = 'transactionStatusSubtitle',
  transactionStatusButton = 'transactionStatusButton',
  heroTitle = 'heroTitle',
  invalidAmount = 'invalidAmount',
  modalCloseButton = 'modalCloseButton',
  cancelTrxBtn = 'cancelTrxBtn',
  transactionToastTitle = 'transactionToastTitle',
  signBtn = 'signBtn',
  closeButton = 'closeButton',
  accessPass = 'accessPass',
  submitButton = 'submitButton',
  keystoreBtn = 'keystoreBtn',
  nextBtn = 'nextBtn',
  confirmBtn = 'confirmBtn',
  typeSubmitSelector = 'button[type="submit"]',
  modalError = '.modal-layout-error',
  tokenSelector = "locator('div').filter({ hasText: /^Select...$/ }).nth(3)"
}

export enum ApiPathEnum {
  devnetGraphdharitriX = 'https://devnet-graph.dharitrix.com/graphql',
  devnetUrl = 'https://devnet-wallet.dharitri.org/'
}

export enum GraphQlProps {
  enabledSwaps = 'enabledSwaps',
  tokenInID = 'tokenInID',
  maintenanceQuery = 'maintenanceQuery'
}

export enum GlobalDataEnum {
  signature = 'bddf92e63149185458521e729955cdc87149f2bf70af7471db4f471824801a2d036503810a20d9e31b86fab910c5c559293b4ee3a8d66d96265a24cac76a9103',
  encodedMsg = '0x54657374206d7367',
  decodedMsg = ' Test msg',
  contractInput = 'drt1qqqqqqqqqqqqqpgqfcm6l6rd42hwhskmk4thlp9kz58npfq50gfqslu5rr',
  contractDevnet = 'drt1qqqqqqqqqqqqqpgqhwsytpz5q57my2ns4t2pml3r22xt9p8fcxuqwwgfju',
  nftImg = 'https://images.app.goo.gl/uCTQH1oUsqmJRd6N9'
}

export enum LoginFilesEnum {
  keystorePath = './playwright-tests/utils/testKeystore1.json',
  keystorePath2 = './playwright-tests/utils/testKeystore2.json',
  pem = './playwright-tests/utils/sovereignTest.pem'
}
