import { MainLayout } from '@/components/Layout'
import {
  ProductForm,
  ProductFormValues,
  productToFormValues,
  useGetProduct,
  useGetProductCategories,
  useUpdateProduct,
} from '@/features/products'
import { Container, Stack, Title } from '@mantine/core'
import { notifications } from '@mantine/notifications'
import { IconCheck } from '@tabler/icons-react'
import { useRouter } from 'next/router'
import { useCallback } from 'react'

const IMAGE_SIZE = 240

interface EditProductProps {
  productId: string
}

export function EditProduct({ productId }: EditProductProps) {
  const router = useRouter()
  const { data: categories, isLoading: isCategoriesLoading } =
    useGetProductCategories(1, 100)
  const { data: product, isLoading: isProductLoading } =
    useGetProduct(productId)

  const allCategories = categories?.data.map((category) => ({
    id: category.id,
    name: category.name,
  }))

  const updateProductMutation = useUpdateProduct({
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
      updateProductMutation.mutate({
        id: productId,
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
    [productId, updateProductMutation]
  )

  const handleCancel = useCallback(() => {
    router.push('/dashboard/products')
  }, [router])

  return (
    <MainLayout>
      <Container size="xs">
        <Stack>
          <Title order={1}>Edit Product</Title>
          {!isProductLoading &&
            !isCategoriesLoading &&
            allCategories !== undefined &&
            product !== undefined && (
              <ProductForm
                imageSize={IMAGE_SIZE}
                allCategories={allCategories}
                initialValues={productToFormValues(product)}
                onSubmit={handleSubmit}
                onCancel={handleCancel}
              />
            )}
        </Stack>
      </Container>
    </MainLayout>
  )
}
