import { getCurrentUser } from '@/features/auth'
import { Session, SessionProvider, SessionStatus } from '@/providers/session'
import { MantineProvider } from '@mantine/core'
import { Notifications } from '@mantine/notifications'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { AxiosError } from 'axios'
import type { AppProps } from 'next/app'
import { useRouter } from 'next/router'
import { useEffect, useState } from 'react'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
      refetchOnWindowFocus: false,
    },
  },
})

export default function App({ Component, pageProps }: AppProps) {
  const router = useRouter()
  const [session, setSession] = useState<Session>()
  const [sessionStatus, setSessionStatus] = useState<SessionStatus>('loading')

  useEffect(() => {
    getCurrentUser()
      .then(({ data: { name, email } }) => {
        if (name !== undefined && email !== undefined) {
          setSession({ user: { name, email } })
          setSessionStatus('authenticated')
        } else {
          throw new Error('username or email is undefined')
        }
      })
      .catch((error: AxiosError) => {
        if (error.response?.status === 401) {
          setSession(undefined)
          setSessionStatus('unauthenticated')
        } else {
          throw new Error(error.message)
        }
      })
  }, [router])

  return (
    <SessionProvider session={session} status={sessionStatus}>
      <MantineProvider withNormalizeCSS withGlobalStyles>
        <QueryClientProvider client={queryClient}>
          <Component {...pageProps} />
          <Notifications position="top-center" />
          <ReactQueryDevtools initialIsOpen={false} />
        </QueryClientProvider>
      </MantineProvider>
    </SessionProvider>
  )
}
