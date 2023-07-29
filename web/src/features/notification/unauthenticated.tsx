import { notifications } from '@mantine/notifications'
import { IconX } from '@tabler/icons-react'
import { NotificationConfig } from './shared'

export const NOTIFY_UNAUTHENTICATED_ERROR_ID = 'unauthorized-error'

export const NOTIFY_UNAUTHENTICATED_ERROR_MESSAGES = {
  AccessToAuthenticatedPage:
    'You must be logged in to access authenticated pages.',
  ExpiredSession: 'Your session has expired. Please log in again.',
  AddToCart: 'You must be logged in to add products to your cart.',
} as const

export function notifyUnauthenticatedError(config: NotificationConfig) {
  notifications.show({
    title: 'Unauthorized Error',
    color: 'red',
    icon: <IconX />,
    withCloseButton: true,
    withBorder: true,
    ...config,
  })
}
