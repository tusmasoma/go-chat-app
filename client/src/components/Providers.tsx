import React, { PropsWithChildren } from 'react';
import { NextUIProvider } from '@nextui-org/react';

export const Providers: React.FC<PropsWithChildren> = ({ children }) => {
  return <NextUIProvider>{children}</NextUIProvider>;
};
