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
import { useSmallerThan } from '@/hooks'
import { Container, Flex, Stack, Title, rem } from '@mantine/core'
import { useRouter } from 'next/router'
import { useCallback } from 'react'
import { OrderSummary } from './OrderSummary'

const DEFAULT_IMAGE_SIZE = 200
const SMALL_IMAGE_SIZE = 96

export function Cart() {
  const router = useRouter()
  const smallerThanSm = useSmallerThan('sm')
  const smallerThanLg = useSmallerThan('lg')
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
      {smallerThanLg ? (
        <Container size="sm">
          <Stack>
            <Title order={1}>Cart</Title>
            <CartProductList
              cartProducts={cart?.products}
              isLoading={isLoading}
              imageSize={smallerThanSm ? SMALL_IMAGE_SIZE : DEFAULT_IMAGE_SIZE}
              onChangeQuantity={handleChangeProductQuantity}
              onDelete={handleDeleteProduct}
            />
            {isLoading && cart === undefined && (
              <OrderSummary cart={cart} isLoading={isLoading} />
            )}
            {cart && cart.products.length > 0 && <OrderSummary cart={cart} />}
          </Stack>
        </Container>
      ) : (
        <Container size="lg">
          <Stack>
            <Title order={1}>Cart</Title>
            <Flex gap="xl">
              <div style={{ flex: '1' }}>
                <CartProductList
                  cartProducts={cart?.products}
                  isLoading={isLoading}
                  imageSize={
                    smallerThanSm ? SMALL_IMAGE_SIZE : DEFAULT_IMAGE_SIZE
                  }
                  onChangeQuantity={handleChangeProductQuantity}
                  onDelete={handleDeleteProduct}
                />
              </div>
              <div style={{ width: rem(400) }}>
                {isLoading && cart === undefined && (
                  <OrderSummary cart={cart} isLoading={isLoading} />
                )}
                {cart && cart.products.length > 0 && (
                  <OrderSummary cart={cart} />
                )}
              </div>
            </Flex>
          </Stack>
        </Container>
      )}
    </MainLayout>
  )
}
