import { notifications } from '@mantine/notifications'
import { IconX } from '@tabler/icons-react'
import { NotificationConfig } from './shared'

export const NOTIFY_UNAUTHORIZED_ERRORS = {
  AccessToAuthenticatedPage: {
    id: 'access-to-authenticated-page-unauthorized-error',
    message: 'You must be logged in to access authenticated pages.',
  },

  ExpiredSession: {
    id: 'expired-session-unauthorized-error',
    message: 'Your session has expired. Please log in again.',
  },

  AddToCart: {
    id: 'add-to-cart-unauthorized-error',
    message: 'You must be logged in to add products to your cart.',
  },
} as const

export function notifyUnauthorizedError({ id, message }: NotificationConfig) {
  notifications.show({
    id,
    message,
    title: 'Unauthorized Error',
    color: 'red',
    icon: <IconX />,
    withCloseButton: true,
    withBorder: true,
  })
}
