import { notifications } from '@mantine/notifications'
import { IconX } from '@tabler/icons-react'
import { NotificationConfig } from './shared'

export const NOTIFY_NOT_IMPLEMENTED_ERRORS = {
  viewMoreProducts: {
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
