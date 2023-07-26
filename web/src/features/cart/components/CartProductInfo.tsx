import { NativeNumberSelect, useForm } from '@/components/Form'
import { Image } from '@/components/Image'
import { Price } from '@/components/Price'
import { zodResolver } from '@hookform/resolvers/zod'
import {
  CloseButton,
  Flex,
  Stack,
  Text,
  clsx,
  createStyles,
  rem,
} from '@mantine/core'
import range from 'lodash/range'
import { ChangeEvent, useCallback, useEffect, useState } from 'react'
import { z } from 'zod'
import { CartProduct } from '../types'

const schema = z.object({
  quantity: z.number().min(1).max(10),
})

const useStyles = createStyles(() => ({
  root: {
    listStyle: 'none',
  },
}))

interface CartProductInfoProps {
  className?: string
  cartProduct: CartProduct
  imageSize: number
  editable?: boolean
  onChangeQuantity?: (id: string, quantity: number) => void
  onDelete?: (id: string) => void
}

export function CartProductInfo({
  className,
  cartProduct,
  imageSize,
  editable = true,
  onChangeQuantity,
  onDelete,
}: CartProductInfoProps) {
  const { classes } = useStyles()

  const [isDeleting, setIsDeleting] = useState(false)

  const [Form, methods] = useForm<{
    quantity: number
  }>({
    resolver: zodResolver(schema),
    defaultValues: {
      quantity: cartProduct.quantity,
    },
  })

  useEffect(() => {
    const { setValue } = methods
    setValue('quantity', cartProduct.quantity)
  }, [methods, cartProduct.quantity])

  const handleDelete = useCallback(() => {
    setIsDeleting(true)
    onDelete?.(cartProduct.id)
  }, [onDelete, cartProduct.id])

  const handleChangeQuantity = useCallback(
    (event: ChangeEvent<HTMLSelectElement>) => {
      const { value } = event.target
      onChangeQuantity?.(cartProduct.id, Number(value))
    },
    [onChangeQuantity, cartProduct.id]
  )

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
              <Form>
                <NativeNumberSelect
                  w={rem(80)}
                  label="Quantity"
                  name="quantity"
                  options={range(1, Math.max(10, cartProduct.quantity) + 1)}
                  onChange={handleChangeQuantity}
                  disabled={isDeleting}
                />
              </Form>
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
