import {
  NOTIFY_UNAUTHENTICATED_ERROR_ID,
  NOTIFY_UNAUTHENTICATED_ERROR_MESSAGES,
  notifyUnauthenticatedError,
} from '@/features/notification/unauthenticated'
import { useSession } from '@/providers/session'
import { useRouter } from 'next/router'
import { useEffect } from 'react'

const REDIRECT_PATH_WHEN_UNAUTHENTICATED = '/'

interface AuthGuardProps {
  children: React.ReactNode
}

export function AuthGuard({ children }: AuthGuardProps) {
  const { status } = useSession()
  const router = useRouter()

  useEffect(() => {
    if (
      status === 'unauthenticated' &&
      router.pathname !== REDIRECT_PATH_WHEN_UNAUTHENTICATED
    ) {
      router.push(REDIRECT_PATH_WHEN_UNAUTHENTICATED)
      notifyUnauthenticatedError({
        id: NOTIFY_UNAUTHENTICATED_ERROR_ID,
        message:
          NOTIFY_UNAUTHENTICATED_ERROR_MESSAGES.AccessToAuthenticatedPage,
      })
    }
  }, [router, status])

  if (status !== 'authenticated') {
    return null
  }

  return <>{children}</>
}
