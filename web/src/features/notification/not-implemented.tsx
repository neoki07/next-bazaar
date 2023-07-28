import { notifications } from '@mantine/notifications'
import { IconX } from '@tabler/icons-react'
import { NotificationConfig } from './shared'

export const NOTIFY_NOT_IMPLEMENTED_ERRORS = {
  AccountSettings: {
    id: 'account-settings-not-implemented-error',
    message: 'Account Settings is not implemented yet',
  },

  ProceedToCheckout: {
    id: 'proceed-to-checkout-not-implemented-error',
    message: 'Proceed to Checkout is not implemented yet',
  },

  ViewMoreProducts: {
    id: 'view-more-products-not-implemented-error',
    message: 'View More Products is not implemented yet',
  },
} as const

export function notifyNotImplementedError({ id, message }: NotificationConfig) {
  notifications.show({
    id,
    message,
    title: 'Not Implemented Error',
    color: 'yellow',
    icon: <IconX />,
    withCloseButton: true,
    withBorder: true,
  })
}
