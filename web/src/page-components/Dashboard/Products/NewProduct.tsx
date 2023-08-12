import { MainLayout } from '@/components/Layout'
import {
  ProductForm,
  ProductFormValues,
  useAddProduct,
  useGetProductCategories,
} from '@/features/products'
import { Container, Stack, Title } from '@mantine/core'
import { useRouter } from 'next/router'
import { useCallback } from 'react'

export function NewProduct() {
  const router = useRouter()
  const { data: categories, isLoading } = useGetProductCategories(1, 100)

  const allCategories = categories?.data.map((category) => ({
    id: category.id,
    name: category.name,
  }))

  const addProductMutation = useAddProduct({
    onSuccess: () => {
      router.push('/dashboard/products')
    },
    onError: (error) => {
      throw new Error(error.message)
    },
  })

  const handleSubmit = useCallback(
    (data: ProductFormValues) => {
      console.error('data:', data)
      addProductMutation.mutate({
        data: {
          name: data.name,
          description: data.description,
          price: data.price.toString(),
          stock_quantity: data.stockQuantity,
          category_id: data.categoryId,
        },
      })
    },
    [addProductMutation]
  )

  return (
    <MainLayout>
      <Container size="xs">
        <Stack>
          <Title order={1}>Add Product</Title>
          {!isLoading && allCategories !== undefined && (
            <ProductForm
              allCategories={allCategories}
              onSubmit={handleSubmit}
            />
          )}
        </Stack>
      </Container>
    </MainLayout>
  )
}
