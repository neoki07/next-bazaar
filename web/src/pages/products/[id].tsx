import { NumberSelect, useForm } from '@/components/Form'
import { MainLayout } from '@/components/Layout'
import { Price, PriceSkeleton } from '@/components/Price'
import { useCartProductsCount } from '@/features/cart'
import { useAddToCart } from '@/features/cart/hooks/useAddToCart'
import { useGetProduct } from '@/features/products'
import { useSession } from '@/providers/session'
import { zodResolver } from '@hookform/resolvers/zod'
import {
  Button,
  Container,
  Flex,
  Image,
  Skeleton,
  Stack,
  Text,
  rem,
} from '@mantine/core'
import { notifications } from '@mantine/notifications'
import { IconX } from '@tabler/icons-react'
import { useRouter } from 'next/router'
import { useCallback } from 'react'
import { z } from 'zod'

function notifyUnauthorizedError() {
  notifications.show({
    id: 'unauthorized-error',
    title: 'Unauthorized Error',
    message: 'You must be logged in to add products to your cart.',
    color: 'red',
    icon: <IconX />,
    withCloseButton: true,
    withBorder: true,
  })
}

interface ProductAreaProps {
  productId: string
}

export function ProductArea({ productId }: ProductAreaProps) {
  const { session } = useSession()
  const { data: product, isLoading } = useGetProduct(productId)
  const { refetch: refetchCartProductsCount } = useCartProductsCount({
    enabled: false,
  })

  const addToCartMutation = useAddToCart({
    onSuccess: () => refetchCartProductsCount(),
    onError: (error) => {
      if (error.response?.status === 401) {
        notifyUnauthorizedError()
      } else {
        throw new Error('Unexpected error')
      }
    },
  })

  const schema = z.object({
    amount: z.number().min(1).max(10),
  })

  const handleSubmit = useCallback(
    (data: z.infer<typeof schema>) => {
      if (session === undefined) {
        notifyUnauthorizedError()
        return
      }

      if (product === undefined) {
        throw new Error('product is undefined')
      }

      addToCartMutation.mutate({
        data: { product_id: product.id, quantity: data.amount },
      })
    },
    [session, product, addToCartMutation]
  )

  const [Form] = useForm<{
    amount: number
  }>({
    resolver: zodResolver(schema),
    defaultValues: {
      amount: 1,
    },
    onSubmit: handleSubmit,
  })

  return (
    <Stack>
      <Skeleton visible={isLoading} width="50%">
        <Text size={28}>{isLoading ? 'dummy' : product?.name}</Text>
      </Skeleton>

      {isLoading ? (
        <div>
          <Skeleton height={12} my={12} />
          <Skeleton height={12} my={12} />
          <Skeleton height={12} my={12} width="50%" />
        </div>
      ) : (
        <Text sx={(theme) => ({ color: theme.colors.gray[7] })}>
          {product?.description}
        </Text>
      )}

      <Flex gap={rem(40)}>
        {isLoading ? (
          <Skeleton>
            {/* TODO: improve the way height is decided */}
            <svg viewBox="0 0 648 641.21" />
          </Skeleton>
        ) : (
          <Image src={product?.imageUrl} alt={product?.name} />
        )}
        <Form>
          <Stack w={rem(240)}>
            {isLoading || product === undefined ? (
              <PriceSkeleton width="50%" />
            ) : (
              <Price price={product.price} />
            )}
            <Skeleton visible={isLoading} width="35%">
              <NumberSelect
                w={rem(80)}
                label="Amount"
                name="amount"
                options={[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}
              />
            </Skeleton>
            <Skeleton visible={isLoading}>
              <Button type="submit" color="dark" fullWidth>
                Add to Cart
              </Button>
            </Skeleton>
          </Stack>
        </Form>
      </Flex>
    </Stack>
  )
}

interface ProductAreaProps {
  productId: string
}

export default function ProductPage() {
  const router = useRouter()
  const { id } = router.query
  if (Array.isArray(id)) {
    throw new Error('id is array:' + JSON.stringify(id))
  }

  return (
    <MainLayout>
      <Container>
        {id !== undefined && <ProductArea productId={id} />}
      </Container>
    </MainLayout>
  )
}
