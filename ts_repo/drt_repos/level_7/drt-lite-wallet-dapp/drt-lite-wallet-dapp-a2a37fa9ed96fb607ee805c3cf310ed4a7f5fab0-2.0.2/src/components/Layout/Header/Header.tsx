import { useMatch } from 'react-router-dom';
import { DrtLink } from 'components/DrtLink';
import { useGetIsLoggedIn } from 'lib';
import { DataTestIdsEnum, RouteNamesEnum } from 'localConstants';
import { routeNames } from 'routes';
import DharitrILogo from '../../../assets/img/dharitri-logo.svg?react';
import { NetworkSwitcher } from '../../NetworkSwitcher';

export const Header = () => {
  const isLoggedIn = useGetIsLoggedIn();
  const isUnlockRoute = Boolean(useMatch(RouteNamesEnum.unlock));

  const ConnectButton = isUnlockRoute ? null : (
    <DrtLink data-testid={DataTestIdsEnum.connectBtn} to={RouteNamesEnum.unlock}>
      Connect
    </DrtLink>
  );

  return (
    <header className='flex flex-row align-center justify-between pl-6 pr-6 pt-6'>
      <DrtLink
        className='flex items-center justify-between'
        to={isLoggedIn ? routeNames.dashboard : routeNames.home}
      >
        <DharitrILogo className='w-full h-6' />
      </DrtLink>

      <nav className='h-full w-full text-sm sm:relative sm:left-auto sm:top-auto sm:flex sm:w-auto sm:flex-row sm:justify-end sm:bg-transparent'>
        <div className='flex justify-end container mx-auto items-center gap-2'>
          <NetworkSwitcher />

          {isLoggedIn ? (
            <DrtLink
              className='inline-block rounded-lg px-3 py-2 text-center hover:no-underline my-0 text-gray-600 hover:bg-slate-100 mx-0'
              data-testid={DataTestIdsEnum.logoutBtn}
              to={routeNames.logout}
            >
              Close
            </DrtLink>
          ) : (
            ConnectButton
          )}
        </div>
      </nav>
    </header>
  );
};
