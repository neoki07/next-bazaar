import { AuthGuard } from '@/features/auth'
import { SessionProvider } from '@/providers/session'
import { PageAuthConfig } from '@/types/page'
import { MantineProvider } from '@mantine/core'
import { Notifications } from '@mantine/notifications'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { NextComponentType } from 'next'
import type { AppProps as NextAppProps } from 'next/app'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
      refetchOnWindowFocus: false,
    },
  },
})

interface AppProps extends NextAppProps {
  Component: NextComponentType & PageAuthConfig
}

export default function App({ Component, pageProps }: AppProps) {
  return (
    <SessionProvider>
      <MantineProvider withNormalizeCSS withGlobalStyles>
        <QueryClientProvider client={queryClient}>
          {Component.auth ? (
            <AuthGuard>
              <Component {...pageProps} />
            </AuthGuard>
          ) : (
            <Component {...pageProps} />
          )}
          <Notifications position="top-center" />
          <ReactQueryDevtools initialIsOpen={false} />
        </QueryClientProvider>
      </MantineProvider>
    </SessionProvider>
  )
}
