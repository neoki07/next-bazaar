import { MainLayout } from '@/components/Layout'
import { Price } from '@/components/Price'
import {
  CartProductList,
  useCart,
  useCartProductsCount,
  useUpdateProductQuantity,
} from '@/features/cart'
import { useDeleteProduct } from '@/features/cart/hooks/useDeleteProduct'
import {
  Button,
  Container,
  Divider,
  Flex,
  Group,
  Stack,
  Text,
  Title,
  createStyles,
  rem,
} from '@mantine/core'
import { useCallback } from 'react'

const IMAGE_SIZE = 200

const useStyles = createStyles((theme) => ({
  container: {
    backgroundColor: theme.colors.gray[1],
    width: rem(400),
    padding: rem(32),
    borderRadius: theme.radius.sm,
  },

  border: {
    opacity: 0.3,
  },

  detailText: {
    fontSize: rem(14),
    color: theme.colors.gray[7],
  },

  detailPriceText: {
    fontSize: rem(14),
    color: theme.colors.gray[8],
    fontWeight: 500,
  },

  orderTotalText: {
    fontSize: rem(16),
    fontWeight: 500,
  },

  orderTotalPriceText: {
    fontSize: rem(16),
    fontWeight: 500,
  },
}))

export default function CartPage() {
  const { classes } = useStyles()

  const { data: cart, isLoading, refetch: refetchCart } = useCart()
  const { refetch: refetchCartProductsCount } = useCartProductsCount({
    enabled: false,
  })

  const updateProductQuantityMutation = useUpdateProductQuantity({
    onSuccess: () => {
      refetchCart()
      refetchCartProductsCount()
    },
  })
  const deleteProductMutation = useDeleteProduct({
    onSuccess: () => {
      refetchCart()
      refetchCartProductsCount()
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
          <Stack className={classes.container} spacing="xl">
            {cart !== undefined && (
              <>
                <Title order={3}>Order summary</Title>
                <div>
                  <Flex align="center">
                    <div style={{ flex: '1' }}>
                      <Text className={classes.detailText}>Subtotal</Text>
                    </div>
                    <Price
                      price={cart.subtotal}
                      className={classes.detailPriceText}
                    />
                  </Flex>

                  <Divider my="sm" />

                  <Flex align="center">
                    <div style={{ flex: '1' }}>
                      <Text className={classes.detailText}>Shipping</Text>
                    </div>
                    <Price
                      price={cart.shipping}
                      className={classes.detailPriceText}
                    />
                  </Flex>

                  <Divider my="sm" />

                  <Flex align="center">
                    <div style={{ flex: '1' }}>
                      <Text className={classes.detailText}>Tax</Text>
                    </div>
                    <Price
                      price={cart.tax}
                      className={classes.detailPriceText}
                    />
                  </Flex>

                  <Divider my="sm" />

                  <Flex align="center">
                    <div style={{ flex: '1' }}>
                      <Text className={classes.orderTotalText}>
                        Order total
                      </Text>
                    </div>
                    <Price
                      price={cart.total}
                      className={classes.orderTotalPriceText}
                    />
                  </Flex>
                </div>
                <Button color="dark" fullWidth>
                  Proceed to Checkout
                </Button>
              </>
            )}
          </Stack>
        </Group>
      </Container>
    </MainLayout>
  )
}

CartPage.auth = true
