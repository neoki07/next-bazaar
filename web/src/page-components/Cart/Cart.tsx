import { MainLayout } from '@/components/Layout'
import {
  CartProductList,
  useCart,
  useCartProductsCount,
  useUpdateProductQuantity,
} from '@/features/cart'
import { useDeleteProduct } from '@/features/cart/hooks/useDeleteProduct'
import {
  NOTIFY_UNAUTHENTICATED_ERROR_ID,
  NOTIFY_UNAUTHENTICATED_ERROR_MESSAGES,
  notifyUnauthenticatedError,
} from '@/features/notification/unauthenticated'
import { Container, Group, Title, rem } from '@mantine/core'
import { useRouter } from 'next/router'
import { useCallback } from 'react'
import { OrderSummary } from './OrderSummary'

const IMAGE_SIZE = 200

export function Cart() {
  const router = useRouter()

  const { data: cart, isLoading, refetch: refetchCart } = useCart()
  const { refetch: refetchCartProductsCount } = useCartProductsCount({
    enabled: false,
  })

  const updateProductQuantityMutation = useUpdateProductQuantity({
    onSuccess: () => {
      refetchCart()
      refetchCartProductsCount()
    },
    onError: (error) => {
      if (error.response?.status === 401) {
        router.push('/')
        notifyUnauthenticatedError({
          id: NOTIFY_UNAUTHENTICATED_ERROR_ID,
          message: NOTIFY_UNAUTHENTICATED_ERROR_MESSAGES.ExpiredSession,
        })
      } else {
        throw new Error('Unexpected error')
      }
    },
  })
  const deleteProductMutation = useDeleteProduct({
    onSuccess: () => {
      refetchCart()
      refetchCartProductsCount()
    },
    onError: (error) => {
      if (error.response?.status === 401) {
        router.push('/')
        notifyUnauthenticatedError({
          id: NOTIFY_UNAUTHENTICATED_ERROR_ID,
          message: NOTIFY_UNAUTHENTICATED_ERROR_MESSAGES.ExpiredSession,
        })
      } else {
        throw new Error('Unexpected error')
      }
    },
  })

  const handleChangeProductQuantity = useCallback(
    (productId: string, quantity: number) => {
      updateProductQuantityMutation.mutate({
        productId,
        data: { quantity },
      })
    },
    [updateProductQuantityMutation]
  )

  const handleDeleteProduct = useCallback(
    (productId: string) => {
      deleteProductMutation.mutate({
        productId,
      })
    },
    [deleteProductMutation]
  )

  return (
    <MainLayout>
      <Container size={1200} miw={800}>
        <Title mb="lg" order={1}>
          Cart
        </Title>
        <Group align="start" spacing={rem(48)}>
          <div style={{ flex: '1' }}>
            <CartProductList
              cartProducts={cart?.products}
              isLoading={isLoading}
              imageSize={IMAGE_SIZE}
              onChangeQuantity={handleChangeProductQuantity}
              onDelete={handleDeleteProduct}
            />
          </div>
          {cart && cart.products.length > 0 && <OrderSummary cart={cart} />}
        </Group>
      </Container>
    </MainLayout>
  )
}
