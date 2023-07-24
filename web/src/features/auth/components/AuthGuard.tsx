import { useSession } from '@/providers/session'
import { notifications } from '@mantine/notifications'
import { IconX } from '@tabler/icons-react'
import { useRouter } from 'next/router'
import { useEffect } from 'react'

const REDIRECT_PATH_WHEN_UNAUTHENTICATED = '/'

function notifyUnauthorizedError() {
  notifications.show({
    id: 'access-to-authenticated-page-unauthorized-error',
    title: 'Unauthorized Error',
    message: 'You must be logged in to access authenticated pages.',
    color: 'red',
    icon: <IconX />,
    withCloseButton: true,
    withBorder: true,
  })
}

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
      notifyUnauthorizedError()
    }
  }, [router, status])

  if (status !== 'authenticated') {
    return null
  }

  return <>{children}</>
}
