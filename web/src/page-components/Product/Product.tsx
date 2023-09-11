import { useForm } from '@/components/Form'
import { NativeNumberSelect } from '@/components/Form/components/NativeNumberSelect'
import { ResponsiveSquareImage } from '@/components/Image'
import { MainLayout } from '@/components/Layout'
import { Price, PriceSkeleton } from '@/components/Price'
import { useCartProductsCount } from '@/features/cart'
import { useAddToCart } from '@/features/cart/hooks/useAddToCart'
import {
  NOTIFY_UNAUTHENTICATED_ERROR_ID,
  NOTIFY_UNAUTHENTICATED_ERROR_MESSAGES,
  notifyUnauthenticatedError,
} from '@/features/notification/unauthenticated'
import { useGetProduct } from '@/features/products'
import { useSmallerThan } from '@/hooks'
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
import range from 'lodash/range'
import { useRouter } from 'next/router'
import { useCallback } from 'react'
import { z } from 'zod'
import { AddedModal } from './AddedModal'

const IMAGE_SIZE = 648

interface ProductAreaProps {
  productId: string
}

export function ProductArea({ productId }: ProductAreaProps) {
  const smallerThanMd = useSmallerThan('sm')
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
        notifyUnauthenticatedError({
          id: NOTIFY_UNAUTHENTICATED_ERROR_ID,
          message: NOTIFY_UNAUTHENTICATED_ERROR_MESSAGES.AddToCart,
        })
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
        notifyUnauthenticatedError({
          id: NOTIFY_UNAUTHENTICATED_ERROR_ID,
          message: NOTIFY_UNAUTHENTICATED_ERROR_MESSAGES.AddToCart,
        })
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
      <Skeleton visible={isLoading} width={isLoading ? '50%' : undefined}>
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

      <Flex direction={smallerThanMd ? 'column' : 'row'} gap="xl">
        <div style={{ flex: 1 }}>
          {isLoading ||
          product === undefined ||
          product.imageUrl === undefined ? (
            <ResponsiveSquareImage isLoading />
          ) : (
            <ResponsiveSquareImage src={product.imageUrl} alt={product.name} />
          )}
        </div>

        <Form>
          <Stack w={rem(240)}>
            {isLoading || product === undefined ? (
              <PriceSkeleton width="50%" size="xl" weight="bold" />
            ) : (
              <Price price={product.price} size="xl" weight="bold" />
            )}
            <Skeleton visible={isLoading} width={isLoading ? '35%' : undefined}>
              <NativeNumberSelect
                w={rem(80)}
                label="Amount"
                name="amount"
                options={range(1, 11)}
              />
            </Skeleton>
            <Skeleton visible={isLoading} width={isLoading ? '35%' : undefined}>
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
