import { NativeNumberSelect, useForm } from '@/components/Form'
import { Image } from '@/components/Image'
import { Price } from '@/components/Price'
import { zodResolver } from '@hookform/resolvers/zod'
import { CloseButton, Flex, Group, Stack, Text, rem } from '@mantine/core'
import range from 'lodash/range'
import { ChangeEvent, useCallback, useEffect, useState } from 'react'
import { z } from 'zod'
import { useCartProducts } from '../hooks/useCartProducts'
import { useCartProductsCount } from '../hooks/useCartProductsCount'
import { useDeleteProduct } from '../hooks/useDeleteProduct'
import { useUpdateProductQuantity } from '../hooks/useUpdateProductQuantity'
import { CartProduct } from '../types'

interface CartProductInfoProps {
  cartProduct: CartProduct
  imageSize: number
}

export function CartProductInfo({
  cartProduct,
  imageSize,
}: CartProductInfoProps) {
  const [isDeleting, setIsDeleting] = useState(false)

  const { refetch: refetchCartProducts } = useCartProducts()
  const { refetch: refetchCartProductsCount } = useCartProductsCount({
    enabled: false,
  })

  const updateProductQuantityMutation = useUpdateProductQuantity({
    onSuccess: () => {
      refetchCartProducts()
      refetchCartProductsCount()
    },
  })
  const deleteProductMutation = useDeleteProduct({
    onSuccess: () => {
      refetchCartProducts()
      refetchCartProductsCount()
    },
  })

  const schema = z.object({
    quantity: z.number().min(1).max(10),
  })

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

  const handleSubmit = useCallback(
    (data: z.infer<typeof schema>) => {
      updateProductQuantityMutation.mutate({
        productId: cartProduct.id,
        data: { quantity: data.quantity },
      })
    },
    [updateProductQuantityMutation, cartProduct.id]
  )

  const handleDelete = useCallback(() => {
    setIsDeleting(true)
    deleteProductMutation.mutate({
      productId: cartProduct.id,
    })
  }, [deleteProductMutation, cartProduct.id])

  const handleChangeQuantity = useCallback(
    (event: ChangeEvent<HTMLSelectElement>) => {
      const { value } = event.target
      handleSubmit({ quantity: Number(value) })
    },
    [handleSubmit]
  )

  return (
    <Group my="sm">
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
          <Form>
            <Stack spacing="xs">
              <Text fz="md">{cartProduct.name}</Text>
              <Price price={cartProduct.price} />
              <NativeNumberSelect
                w={rem(80)}
                label="Quantity"
                name="quantity"
                options={range(1, Math.max(10, cartProduct.quantity) + 1)}
                onChange={handleChangeQuantity}
                disabled={isDeleting}
              />
            </Stack>
          </Form>
          <CloseButton
            aria-label="Remove product"
            onClick={handleDelete}
            disabled={isDeleting}
          />
        </Flex>
      </Stack>
    </Group>
  )
}
