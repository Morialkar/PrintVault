'use client';

import { Provider, createClient } from '@printvault/api';
import React from 'react';

export const FuseProvider = (props: any) => {
  const [client, ssr] = React.useMemo(() => {
    const { client, ssr } = createClient({
      url: 'https://ideal-space-couscous-r64p46qxrx2xx67-4000.app.github.dev/graphql',
      // This is used during SSR to know when the data finishes loading
      suspense: true,
    });

    return [client, ssr];
  }, []);

  return (
    <Provider client={client} ssr={ssr}>
      <React.Suspense>{props.children}</React.Suspense>
    </Provider>
  );
};
