import { MainLayout } from '@/components/Layout'
import {
  ProductForm,
  ProductFormValues,
  useAddProduct,
  useGetProductCategories,
} from '@/features/products'
import { useSmallerThan } from '@/hooks'
import { Container, Stack, Title, rem } from '@mantine/core'
import { notifications } from '@mantine/notifications'
import { IconCheck } from '@tabler/icons-react'
import { useRouter } from 'next/router'
import { useCallback } from 'react'

const DEFAULT_IMAGE_SIZE = 300
const SMALL_IMAGE_SIZE = 200

export function NewProduct() {
  const router = useRouter()
  const smallerThan440px = useSmallerThan(rem(440))
  const { data: categories, isLoading } = useGetProductCategories(1, 100)

  const allCategories = categories?.data.map((category) => ({
    id: category.id,
    name: category.name,
  }))

  const addProductMutation = useAddProduct({
    onSuccess: () => {
      router.push('/dashboard/products')

      notifications.show({
        id: 'product-saved-successfully',
        title: 'Product Saved Successfully',
        message: null,
        color: 'teal',
        icon: <IconCheck />,
        withCloseButton: true,
        withBorder: true,
      })
    },
    onError: (error) => {
      throw new Error(error.message)
    },
  })

  const handleSubmit = useCallback(
    (data: ProductFormValues) => {
      addProductMutation.mutate({
        data: {
          name: data.name,
          description: data.description,
          price: data.price.toString(),
          stock_quantity: data.stockQuantity,
          category_id: data.categoryId,
          image_url: data.imageUrl,
        },
      })
    },
    [addProductMutation]
  )

  const handleCancel = useCallback(() => {
    router.push('/dashboard/products')
  }, [router])

  return (
    <MainLayout>
      <Container size="xs">
        <Stack>
          <Title order={1}>Add Product</Title>
          {!isLoading && allCategories !== undefined && (
            <ProductForm
              imageSize={
                smallerThan440px ? SMALL_IMAGE_SIZE : DEFAULT_IMAGE_SIZE
              }
              allCategories={allCategories}
              onSubmit={handleSubmit}
              onCancel={handleCancel}
            />
          )}
        </Stack>
      </Container>
    </MainLayout>
  )
}
