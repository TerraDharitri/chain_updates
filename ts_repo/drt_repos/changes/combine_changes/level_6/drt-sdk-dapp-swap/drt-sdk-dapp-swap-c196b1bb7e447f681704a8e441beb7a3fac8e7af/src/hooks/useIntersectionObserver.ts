import { Dispatch, SetStateAction, useEffect } from 'react';
import { TokensPaginationType, UserDcdtType } from 'types';

interface UseIntersectionObserverType {
  hasMore: boolean;
  observerId: string;
  isLoading: boolean;
  currentCursor: string;
  tokens: UserDcdtType[];
  loadedCursors: Set<string>;
  setPagination: Dispatch<SetStateAction<TokensPaginationType>>;
}

export const useIntersectionObserver = ({
  tokens,
  hasMore,
  isLoading,
  observerId,
  loadedCursors,
  currentCursor,
  setPagination
}: UseIntersectionObserverType) => {
  useEffect(() => {
    if (!observerId) return;
    const element = document.getElementById(observerId);
    if (!element) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore && !isLoading) {
          // Check if the current cursor is already loaded
          if (currentCursor && !loadedCursors.has(currentCursor)) {
            const newPagination = {
              first: 20,
              after: currentCursor
            };

            setPagination(newPagination);
          }
        }
      },
      { threshold: 1.0 }
    );

    observer.observe(element);

    return () => {
      observer.disconnect();
    };
  }, [tokens, observerId]);
};
