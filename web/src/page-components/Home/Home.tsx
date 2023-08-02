import { MainLayout } from '@/components/Layout'
import { Product, useGetProductCategories } from '@/features/products'
import { Stack, rem } from '@mantine/core'
import { range } from 'lodash'
import { useCallback } from 'react'
import { CategorySection } from './CategorySection'
import { CategorySectionSkeleton } from './CategorySectionSkeleton'

const IMAGE_SIZE = 260
const CATEGORY_COUNT_ON_LOAD = 4
const PRODUCT_COUNT_PER_CATEGORY = 8

export function Home() {
  const { data: categories, isLoading } = useGetProductCategories(1, 100)
  const getProductLink = useCallback(
    (product: Product) => `/products/${product.id}`,
    []
  )

  return (
    <MainLayout>
      <Stack spacing={rem(40)}>
        {isLoading ? (
          <>
            {range(CATEGORY_COUNT_ON_LOAD).map((index) => (
              <CategorySectionSkeleton
                key={index}
                imageSize={IMAGE_SIZE}
                productCount={PRODUCT_COUNT_PER_CATEGORY}
              />
            ))}
          </>
        ) : (
          <>
            {categories?.data.map((category) => (
              <CategorySection
                key={category.id}
                category={category}
                getProductLink={getProductLink}
                imageSize={IMAGE_SIZE}
                productCount={PRODUCT_COUNT_PER_CATEGORY}
              />
            ))}
          </>
        )}
      </Stack>
    </MainLayout>
  )
}
