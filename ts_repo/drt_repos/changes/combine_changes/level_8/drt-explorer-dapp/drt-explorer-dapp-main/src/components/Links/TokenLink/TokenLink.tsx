import classNames from 'classnames';
import { useSelector } from 'react-redux';

import { NATIVE_TOKEN_IDENTIFIER } from 'appConstants';
import { NativeTokenSymbol, NetworkLink, Overlay } from 'components';
import { urlBuilder, isRewaToken, isProof } from 'helpers';
import { activeNetworkSelector } from 'redux/selectors';
import { TokenType, TokenTypeEnum } from 'types';

export const TokenLink = ({ token }: { token: TokenType }) => {
  const isNftProof = isProof(token);

  const { rewaLabel = '' } = useSelector(activeNetworkSelector);

  const identifierArray = token.identifier ? token.identifier.split('-') : [];
  if (identifierArray.length > 2 && token.type === TokenTypeEnum.MetaDCDT) {
    identifierArray.pop();
  }
  const detailsIdentifier = identifierArray.join('-');

  const metaDetails = isNftProof
    ? urlBuilder.proofDetails(token.identifier)
    : urlBuilder.tokenMetaDcdtDetails(detailsIdentifier);

  const networkLink =
    token.type === TokenTypeEnum.MetaDCDT
      ? metaDetails
      : urlBuilder.tokenDetails(token.identifier);

  if (token.identifier === NATIVE_TOKEN_IDENTIFIER) {
    const isCustomIcon = !isRewaToken(rewaLabel);
    const defaultCoinLink = `/${rewaLabel.toLowerCase()}`;

    return (
      <NetworkLink to={defaultCoinLink} className='d-flex text-truncate'>
        <span className='fam'></span>
        <NativeTokenSymbol
          className={classNames('sym', { custom: isCustomIcon })}
        />
        <sup className='suf opc'></sup>
      </NetworkLink>
    );
  }

  const TokenComponent = () => (
    <span className='d-flex align-items-center gap-1 text-truncate'>
      {(token.assets?.svgUrl || token.assets?.pngUrl) && (
        <img
          src={token.assets?.svgUrl || token.assets?.pngUrl}
          alt={token.identifier}
          className='side-icon'
        />
      )}
      <div className='text-truncate'>{token?.ticker ?? token.name}</div>
    </span>
  );

  return (
    <NetworkLink
      to={networkLink}
      className={`d-flex text-truncate ${
        token.assets?.svgUrl || token.assets?.pngUrl ? 'side-link' : ''
      }`}
    >
      <div className='d-flex align-items-center symbol text-truncate'>
        {token?.assets ? (
          <>
            {token.type === TokenTypeEnum.MetaDCDT &&
            detailsIdentifier !== token.identifier ? (
              <Overlay title={token.identifier} truncate>
                <TokenComponent />
              </Overlay>
            ) : (
              <TokenComponent />
            )}
          </>
        ) : (
          <span className='text-truncate'>{token.identifier}</span>
        )}
      </div>
    </NetworkLink>
  );
};
