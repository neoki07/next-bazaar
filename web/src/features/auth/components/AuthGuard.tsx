import {
  NOTIFY_UNAUTHORIZED_ERRORS,
  notifyUnauthorizedError,
} from '@/features/notification/unauthorized'
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
      notifyUnauthorizedError(
        NOTIFY_UNAUTHORIZED_ERRORS.accessToAuthenticatedPage
      )
    }
  }, [router, status])

  if (status !== 'authenticated') {
    return null
  }

  return <>{children}</>
}
