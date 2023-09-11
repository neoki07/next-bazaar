import { MainLayout } from '@/components/Layout'
import { Product, useGetProductCategories } from '@/features/products'
import { Container, Stack, rem } from '@mantine/core'
import { range } from 'lodash'
import { useCallback } from 'react'
import { CategorySection } from './CategorySection'
import { CategorySectionSkeleton } from './CategorySectionSkeleton'

const CATEGORY_COUNT_ON_LOAD = 4
const PRODUCT_COUNT_PER_CATEGORY = 12

export function Home() {
  const { data: categories, isLoading } = useGetProductCategories(1, 100)
  const getProductLink = useCallback(
    (product: Product) => `/products/${product.id}`,
    []
  )

  return (
    <MainLayout>
      <Container size="lg">
        <Stack spacing={rem(48)}>
          {isLoading ? (
            <>
              {range(CATEGORY_COUNT_ON_LOAD).map((index) => (
                <div key={index} style={{ width: '100%' }}>
                  <CategorySectionSkeleton
                    key={index}
                    productCount={PRODUCT_COUNT_PER_CATEGORY}
                  />
                </div>
              ))}
            </>
          ) : (
            <>
              {categories?.data.map((category) => (
                <div key={category.id} style={{ width: '100%' }}>
                  <CategorySection
                    category={category}
                    getProductLink={getProductLink}
                    productCount={PRODUCT_COUNT_PER_CATEGORY}
                  />
                </div>
              ))}
            </>
          )}
        </Stack>
      </Container>
    </MainLayout>
  )
}
