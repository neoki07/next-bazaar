import { MainLayout } from '@/components/Layout'
import { Product, useGetProductCategories } from '@/features/products'
import { Stack, rem } from '@mantine/core'
import { useCallback } from 'react'
import { CategorySection } from './CategorySection'

const IMAGE_SIZE = 260

export function Home() {
  const { data: categories, isLoading } = useGetProductCategories(1, 100)
  const getProductLink = useCallback(
    (product: Product) => `/products/${product.id}`,
    []
  )

  if (isLoading) {
    return null
  }

  return (
    <MainLayout>
      <Stack spacing={rem(40)}>
        {categories?.data.map((category) => (
          <CategorySection
            key={category.id}
            category={category}
            getProductLink={getProductLink}
            imageSize={IMAGE_SIZE}
          />
        ))}
      </Stack>
    </MainLayout>
  )
}
