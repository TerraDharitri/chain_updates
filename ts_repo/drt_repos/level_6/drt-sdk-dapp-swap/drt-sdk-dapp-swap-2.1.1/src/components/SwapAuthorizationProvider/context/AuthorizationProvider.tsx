import React, { ReactNode } from 'react';
import {
  AuthorizationContextType,
  AuthorizationContext
} from './AuthorizationContext';

export const AuthorizationProvider = ({
  value,
  children
}: {
  value: AuthorizationContextType;
  children?: ReactNode;
}) => {
  return (
    <AuthorizationContext.Provider value={value}>
      {children}
    </AuthorizationContext.Provider>
  );
};
