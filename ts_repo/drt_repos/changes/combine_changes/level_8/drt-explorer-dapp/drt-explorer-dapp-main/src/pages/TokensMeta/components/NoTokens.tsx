import { PageState } from 'components';
import { faCoins } from 'icons/regular';

export const NoTokens = () => {
  return (
    <PageState
      icon={faCoins}
      title='No Meta-DCDT Tokens'
      className='py-spacer my-auto'
    />
  );
};
