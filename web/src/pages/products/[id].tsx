import { NumberSelect, useForm } from '@/components/Form'
import { MainLayout } from '@/components/Layout'
import { Price } from '@/components/Price'
import { useCartProductsCount } from '@/features/cart'
import { useAddToCart } from '@/features/cart/hooks/useAddToCart'
import { useGetProduct } from '@/features/products'
import { useSession } from '@/providers/session'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button, Container, Flex, Image, Stack, Text, rem } from '@mantine/core'
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

  const [Form, methods] = useForm<{
    amount: number
  }>({
    resolver: zodResolver(schema),
    defaultValues: {
      amount: 1,
    },
    onSubmit: handleSubmit,
  })

  if (isLoading) {
    return <div>Loading...</div>
  }

  if (product === undefined) {
    throw new Error('product is undefined')
  }

  return (
    <Stack>
      <Text size={28}>{product.name}</Text>
      <Text sx={(theme) => ({ color: theme.colors.gray[7] })}>
        {product.description}
      </Text>
      <Flex gap={rem(40)}>
        <Image src={product.imageUrl} alt={product.name} />
        <Form>
          <Stack w={rem(240)}>
            <Price price={product.price} />
            <NumberSelect
              w={rem(80)}
              label="Amount"
              name="amount"
              options={[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}
            />
            <Button type="submit" color="dark" fullWidth>
              Add to Cart
            </Button>
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
