import { newE2EPage } from '@stencil/core/testing';
import { DataTestIdsEnum } from 'constants/dataTestIds.enum';

const tag = 'ledger-connect-panel';

describe('ledger-connect-panel', () => {
  it(`renders ${tag}`, async () => {
    const page = await newE2EPage();

    await page.setContent(`<${tag}></${tag}>`);
    const element = await page.find(tag);
    expect(element).toHaveClass('hydrated');
  });

  it('renders changes according to data', async () => {
    const page = await newE2EPage();

    await page.setContent(`<${tag}></${tag}>`);
    const element = await page.find(`${tag} >>> div`);
    expect(element.textContent).toContain(`Unlock your device`);

    const component = await page.find(tag);
    component.setProperty('data', { connectScreenData: { error: 'Unable to find device' } });
    await page.waitForChanges();
    const error = await page.find(`${tag} >>> ledger-connect-screen >>> div`);
    expect(error.textContent).toContain(`Unable to find device`);

    const loadingUpdate = { connectScreenData: {}, accountScreenData: { accounts: [], startIndex: 0, addressesPerPage: 10, isLoading: true }, confirmScreenData: null };
    component.setProperty('data', loadingUpdate);
    const loading = await page.find(`${tag} >>> [data-testid=${DataTestIdsEnum.ledgerLoading}]`);
    expect(loading).toBeDefined();

    const accountsScreenFetchingUpdate = {
      connectScreenData: {},
      accountScreenData: {
        accounts: [
          { address: 'drt18ykcupcm55vfjzy40tjqck260ktk3shhh7fh5g5exmjxqvm3v7jsupsmga', balance: '...', index: 0 },
          { address: 'drt1894027ex33q3x0nx52w0w0dtlrtvmdpm7jhg8lwyta5s33wp4cps0z4402', balance: '...', index: 1 },
          { address: 'drt1sf95q8xf0y7q237daxwv0w9tuc56krruw5t63ju8e3g7sfcfgrasf3kls3', balance: '...', index: 2 },
          { address: 'drt1cphyth5ud3htxcnpuuttxh0tr0mt3ja8pgyvm7kgf06suzqud9hqtq4jyg', balance: '...', index: 3 },
          { address: 'drt1vm7ds77jq5t424gd56336a57sh84lh302lqjvrc6quserg59vr7q47zmkk', balance: '...', index: 4 },
          { address: 'drt1xt4fu2q3wdzy47aktcc2907qqqtfm4xcffl4tf2zyyfskmhth9uskm6ajl', balance: '...', index: 5 },
          { address: 'drt19mgjxltf6d6qt5saw2f787aypdpj59al95m2hkkanqamchf5dzzqm0r6vx', balance: '...', index: 6 },
          { address: 'drt19hasmrswh6khqspse9ck35vf0eykpkjyq2y0uma8ur0pfmlqqersdw99vt', balance: '...', index: 7 },
          { address: 'drt186l9rguljypvm0lxdhqjg62h3p8yg4ugry69vkhw2j4xsz0ch7kswh7dkf', balance: '...', index: 8 },
          { address: 'drt1ydfcnm4u7syq03vr959ljk5cqfld0vsjgfgejqjshzt6nktnddhsafexqe', balance: '...', index: 9 },
        ],
      },
      confirmScreenData: null,
    };
    component.setProperty('data', accountsScreenFetchingUpdate);
    await page.waitForChanges();
    const fetchingList = await page.find(`${tag} >>> [data-testid=${DataTestIdsEnum.addressTableContainer}]`);
    expect(fetchingList).toBeDefined();

    const balancesUpdate = {
      connectScreenData: {},
      accountScreenData: {
        accounts: [
          { address: 'drt18ykcupcm55vfjzy40tjqck260ktk3shhh7fh5g5exmjxqvm3v7jsupsmga', balance: '8998313560000000000', index: 0 },
          { address: 'drt1894027ex33q3x0nx52w0w0dtlrtvmdpm7jhg8lwyta5s33wp4cps0z4402', balance: '0', index: 1 },
          { address: 'drt1sf95q8xf0y7q237daxwv0w9tuc56krruw5t63ju8e3g7sfcfgrasf3kls3', balance: '139289100230000000', index: 2 },
          { address: 'drt1cphyth5ud3htxcnpuuttxh0tr0mt3ja8pgyvm7kgf06suzqud9hqtq4jyg', balance: '0', index: 3 },
          { address: 'drt1vm7ds77jq5t424gd56336a57sh84lh302lqjvrc6quserg59vr7q47zmkk', balance: '0', index: 4 },
          { address: 'drt1xt4fu2q3wdzy47aktcc2907qqqtfm4xcffl4tf2zyyfskmhth9uskm6ajl', balance: '32299650000000000000', index: 5 },
          { address: 'drt19mgjxltf6d6qt5saw2f787aypdpj59al95m2hkkanqamchf5dzzqm0r6vx', balance: '0', index: 6 },
          { address: 'drt19hasmrswh6khqspse9ck35vf0eykpkjyq2y0uma8ur0pfmlqqersdw99vt', balance: '0', index: 7 },
          { address: 'drt186l9rguljypvm0lxdhqjg62h3p8yg4ugry69vkhw2j4xsz0ch7kswh7dkf', balance: '0', index: 8 },
          { address: 'drt1ydfcnm4u7syq03vr959ljk5cqfld0vsjgfgejqjshzt6nktnddhsafexqe', balance: '0', index: 9 },
        ],
        startIndex: 10,
        addressesPerPage: 10,
        isLoading: false,
      },
      confirmScreenData: null,
    };
    component.setProperty('data', balancesUpdate);
    await page.waitForChanges();
    const addressTableContainer = await page.find(`${tag} >>> [data-testid=${DataTestIdsEnum.addressTableContainer}]`);
    expect(addressTableContainer.textContent).toContain('89983');

    const connectScreenUpdate = {
      connectScreenData: {},
      accountScreenData: null,
      confirmScreenData: {
        selectedAddress: 'drt18ykcupcm55vfjzy40tjqck260ktk3shhh7fh5g5exmjxqvm3v7jsupsmga',
        data: 'https://localhost:3001 for more than one day.',
        confirmAddressText: 'For security, please confirm your address:',
        authText: 'and authorize:',
        areShownText: '',
      },
    };
    component.setProperty('data', connectScreenUpdate);
    await page.waitForChanges();
    const confirmScreen = await page.find(`${tag} >>> [data-testid=${DataTestIdsEnum.ledgerConfirmAddress}]`);
    expect(confirmScreen.textContent).toContain('please confirm your address');
  });
});
