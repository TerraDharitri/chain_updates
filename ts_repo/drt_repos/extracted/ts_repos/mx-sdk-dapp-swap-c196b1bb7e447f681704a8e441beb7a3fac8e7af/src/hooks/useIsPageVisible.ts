import { useCallback, useEffect, useState } from 'react';

const getIsPageVisible = () => {
  if (typeof document === undefined) {
    return true;
  }

  return document.visibilityState === 'visible';
};

export const useIsPageVisible = () => {
  const [isVisible, setIsVisible] = useState(getIsPageVisible());

  const handleVisbilityChange = useCallback(() => {
    setIsVisible(getIsPageVisible());
  }, [setIsVisible]);

  useEffect(() => {
    if (typeof document === undefined) {
      return;
    }

    document.addEventListener('visibilitychange', handleVisbilityChange);
    return () =>
      document.removeEventListener('visibilitychange', handleVisbilityChange);
  }, [handleVisbilityChange]);

  return isVisible;
};
