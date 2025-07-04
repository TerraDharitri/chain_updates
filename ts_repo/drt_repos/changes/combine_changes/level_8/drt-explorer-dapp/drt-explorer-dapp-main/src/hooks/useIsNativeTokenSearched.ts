import { useSelector } from 'react-redux';

import { BRAND_NAME, NATIVE_TOKEN_IDENTIFIER } from 'appConstants';
import { useGetSearch } from 'hooks';
import { activeNetworkSelector } from 'redux/selectors';

export const useIsNativeTokenSearched = () => {
  const { rewaLabel } = useSelector(activeNetworkSelector);
  const { search } = useGetSearch();

  const isNativeTokenSearched = Boolean(
    search &&
      [
        'rewa',
        'numbat',
        'dharitri',
        BRAND_NAME.toLowerCase(),
        (rewaLabel ?? '').toLowerCase(),
        NATIVE_TOKEN_IDENTIFIER.toLowerCase()
      ].includes(search.toLowerCase().trim())
  );

  return isNativeTokenSearched;
};
