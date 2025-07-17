import {
  ApiEndpointsEnum,
  ApiMethodsEnum,
  AssertionEnum,
  RoutesEnum
} from '../../constants/enums';

describe('Tokens', () => {
  beforeEach(() => {
    cy.accesPage(RoutesEnum.tokens);
  });
  it('should display de header elements', () => {
    cy.checkHeaderElements('Tokens');
  });

  it('should properly display the table', () => {
    const tableHead = [
      'Token',
      'Name',
      'Price',
      'Circulating Supply',
      'Market Cap',
      'Holders',
      'Transactions'
    ];
    cy.checkTableHead(tableHead);
  });

  it('should properly change the table page', () => {
    cy.paginationHandler(RoutesEnum.tokens);
  });

  it('should acces the token details page', () => {
    cy.viewport(1000, 3000);
    cy.contains('a', 'WREWA').click();
    cy.getSelector('title').should(
      AssertionEnum.contain,
      'WrappedREWA (WREWA) Token'
    );
  });
  it('should properly search tokens', () => {
    cy.get('header').invoke('css', {
      display: 'none'
    });
    cy.getSelector('tokensSearch').type('$').type('{enter}');
    cy.contains('No Tokens');
    cy.getSelector('tokensSearch').clear().type('REWA').type('{enter}');
    cy.checkUrl('/tokens?search=REWA');
    cy.get('tr').should(AssertionEnum.contain, 'xREWA');
  });
});
