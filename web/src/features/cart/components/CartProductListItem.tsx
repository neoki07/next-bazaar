import { Image } from '@/components/Image'
import { Price } from '@/components/Price'
import {
  CloseButton,
  Flex,
  Stack,
  Text,
  clsx,
  createStyles,
} from '@mantine/core'
import range from 'lodash/range'
import { useCallback, useState } from 'react'
import { CartProduct } from '../types'
import { QuantitySelectForm } from './QuantitySelectForm'

const useStyles = createStyles(() => ({
  root: {
    listStyle: 'none',
  },
}))

interface CartProductListItemProps {
  className?: string
  cartProduct: CartProduct
  imageSize: number
  editable?: boolean
  onChangeQuantity?: (id: string, quantity: number) => void
  onDelete?: (id: string) => void
}

export function CartProductListItem({
  className,
  cartProduct,
  imageSize,
  editable = true,
  onChangeQuantity,
  onDelete,
}: CartProductListItemProps) {
  const { classes } = useStyles()

  const [isDeleting, setIsDeleting] = useState(false)

  const handleDelete = useCallback(() => {
    setIsDeleting(true)
    onDelete?.(cartProduct.id)
  }, [onDelete, cartProduct.id])

  return (
    <li className={clsx(classes.root, className)}>
      <Stack spacing={4}>
        <Flex gap="xs">
          {cartProduct.imageUrl !== undefined && (
            <Image
              src={cartProduct.imageUrl}
              alt={cartProduct.name}
              width={imageSize}
              height={imageSize}
            />
          )}
          <Stack spacing="xs">
            <Text fz="md">{cartProduct.name}</Text>
            <Price price={cartProduct.price} size="xl" weight="bold" />
            {editable ? (
              <QuantitySelectForm
                cartProduct={cartProduct}
                options={range(1, Math.max(10, cartProduct.quantity) + 1)}
                onChange={onChangeQuantity}
                disabled={isDeleting}
              />
            ) : (
              <Text>Quantity: {cartProduct.quantity}</Text>
            )}
          </Stack>
          {editable && (
            <CloseButton
              aria-label="Remove product"
              onClick={handleDelete}
              disabled={isDeleting}
            />
          )}
        </Flex>
      </Stack>
    </li>
  )
}
