import { Price } from '@/components/Price'
import { Cart } from '@/features/cart'
import {
  NOTIFY_NOT_IMPLEMENTED_ERRORS,
  notifyNotImplementedError,
} from '@/features/notification/not-implemented'
import {
  Button,
  Divider,
  Flex,
  Skeleton,
  Stack,
  Text,
  Title,
  createStyles,
  rem,
} from '@mantine/core'

const useStyles = createStyles((theme) => ({
  container: {
    backgroundColor: theme.colors.gray[1],
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

interface LoadingOrderSummaryProps {
  cart?: Cart
  isLoading: true
}

interface NonLoadingOrderSummaryProps {
  cart: Cart
  isLoading?: false
}

type OrderSummaryProps = LoadingOrderSummaryProps | NonLoadingOrderSummaryProps

export function OrderSummary({ cart, isLoading }: OrderSummaryProps) {
  const { classes } = useStyles()

  return (
    <Skeleton visible={!!isLoading}>
      <Stack className={classes.container} spacing="xl">
        <Title order={3}>Order summary</Title>
        <div>
          <Flex align="center">
            <div style={{ flex: '1' }}>
              <Text className={classes.detailText}>Subtotal</Text>
            </div>
            {!isLoading && (
              <Price
                price={cart.subtotal}
                className={classes.detailPriceText}
              />
            )}
          </Flex>

          <Divider my="sm" />

          <Flex align="center">
            <div style={{ flex: '1' }}>
              <Text className={classes.detailText}>Shipping</Text>
            </div>
            {!isLoading && (
              <Price
                price={cart.shipping}
                className={classes.detailPriceText}
              />
            )}
          </Flex>

          <Divider my="sm" />

          <Flex align="center">
            <div style={{ flex: '1' }}>
              <Text className={classes.detailText}>Tax</Text>
            </div>
            {!isLoading && (
              <Price price={cart.tax} className={classes.detailPriceText} />
            )}
          </Flex>

          <Divider my="sm" />

          <Flex align="center">
            <div style={{ flex: '1' }}>
              <Text className={classes.orderTotalText}>Order total</Text>
            </div>
            {!isLoading && (
              <Price
                price={cart.total}
                className={classes.orderTotalPriceText}
              />
            )}
          </Flex>
        </div>
        <Button
          color="dark"
          fullWidth
          onClick={() =>
            notifyNotImplementedError(
              NOTIFY_NOT_IMPLEMENTED_ERRORS.ProceedToCheckout
            )
          }
        >
          Proceed to Checkout
        </Button>
      </Stack>
    </Skeleton>
  )
}
