import { useForm } from '@/components/Form'
import { NativeNumberSelect } from '@/components/Form/components/NativeNumberSelect'
import { Image } from '@/components/Image'
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
  Skeleton,
  Stack,
  Text,
  rem,
} from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { notifications } from '@mantine/notifications'
import { IconX } from '@tabler/icons-react'
import range from 'lodash/range'
import { useRouter } from 'next/router'
import { useCallback } from 'react'
import { z } from 'zod'
import { AddedModal } from './AddedModal'

const IMAGE_SIZE = 648

function notifyUnauthorizedError() {
  notifications.show({
    id: 'add-cart-unauthorized-error',
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
  const [openedModal, { open: openModal, close: closeModal }] =
    useDisclosure(false)
  const { session } = useSession()
  const { data: product, isLoading } = useGetProduct(productId)
  const { refetch: refetchCartProductsCount } = useCartProductsCount({
    enabled: false,
  })

  const addToCartMutation = useAddToCart({
    onSuccess: () => {
      refetchCartProductsCount()
      openModal()
    },
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
        {isLoading ||
        product === undefined ||
        product.imageUrl === undefined ? (
          <Image isLoading alt="" width={IMAGE_SIZE} height={IMAGE_SIZE} />
        ) : (
          <Image
            src={product.imageUrl}
            alt={product.name}
            width={IMAGE_SIZE}
            height={IMAGE_SIZE}
          />
        )}

        <Form>
          <Stack w={rem(240)}>
            {isLoading || product === undefined ? (
              <PriceSkeleton width="50%" size="xl" weight="bold" />
            ) : (
              <Price price={product.price} size="xl" weight="bold" />
            )}
            <Skeleton visible={isLoading} width="35%">
              <NativeNumberSelect
                w={rem(80)}
                label="Amount"
                name="amount"
                options={range(1, 11)}
              />
            </Skeleton>
            <Skeleton visible={isLoading}>
              <Button type="submit" color="dark" fullWidth>
                Add to Cart
              </Button>
              <AddedModal opened={openedModal} onClose={closeModal} />
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

export function Product() {
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
