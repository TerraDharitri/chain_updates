import { SVGProps } from 'react';
import { useSelector } from 'react-redux';

import { ReactComponent as RewaLogo } from 'assets/img/tokens/rewa-logo.svg';
import { ReactComponent as SpcLogo } from 'assets/img/tokens/spc-logo.svg';
import { isRewaToken } from 'helpers';
import { activeNetworkSelector } from 'redux/selectors';

// temporary?
export const NativeTokenLogo = (props: SVGProps<SVGSVGElement>) => {
  const { rewaLabel } = useSelector(activeNetworkSelector);

  if (isRewaToken(rewaLabel)) {
    return <RewaLogo {...props} />;
  }

  switch (rewaLabel?.toLowerCase()) {
    case 'spc':
      return <SpcLogo {...props} />;
    default:
      return null;
  }
};
