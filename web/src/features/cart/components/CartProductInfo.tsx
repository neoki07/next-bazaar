import { NumberSelect, useForm } from '@/components/Form'
import { Price } from '@/components/Price'
import { zodResolver } from '@hookform/resolvers/zod'
import {
  CloseButton,
  Flex,
  Group,
  Image,
  Stack,
  Text,
  rem,
} from '@mantine/core'
import { useCallback, useEffect, useState } from 'react'
import { z } from 'zod'
import { useCartProducts } from '../hooks/useCartProducts'
import { useCartProductsCount } from '../hooks/useCartProductsCount'
import { useDeleteProduct } from '../hooks/useDeleteProduct'
import { useUpdateProductQuantity } from '../hooks/useUpdateProductQuantity'
import { CartProduct } from '../types'

interface CartProductInfoProps {
  cartProduct: CartProduct
}

export function CartProductInfo({ cartProduct }: CartProductInfoProps) {
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
        data: { product_id: cartProduct.id, quantity: data.quantity },
      })
    },
    [updateProductQuantityMutation, cartProduct.id]
  )

  const handleDelete = useCallback(() => {
    setIsDeleting(true)
    deleteProductMutation.mutate({
      data: { product_id: cartProduct.id },
    })
  }, [deleteProductMutation, cartProduct.id])

  const handleChangeQuantity = useCallback(
    (value: string | null) => {
      handleSubmit({ quantity: Number(value) })
    },
    [handleSubmit]
  )

  return (
    <Group my="sm">
      <Stack spacing={4}>
        <Flex gap="xs">
          <Image
            src={cartProduct.imageUrl}
            alt={cartProduct.name}
            width={192}
            height={192}
          />
          <Form>
            <Text fz="md">{cartProduct.name}</Text>
            <Stack spacing={0}>
              <Price price={cartProduct.price} />
              <NumberSelect
                w={rem(80)}
                label="Quantity"
                name="quantity"
                options={[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}
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
