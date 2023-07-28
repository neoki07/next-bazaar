import { notifications } from '@mantine/notifications'
import { IconX } from '@tabler/icons-react'

export const NOT_IMPLEMENTED_ERROR_IDS = {
  viewMoreProducts: 'view-more-products-not-implemented-error',
} as const

export function notifyNotImplementedError(id: string, message: string) {
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
