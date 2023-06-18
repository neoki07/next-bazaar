import { MainLayout } from '@/components/Layout'
import {
  Product,
  ProductCard,
  ProductCardSkeleton,
  useGetProducts,
} from '@/features/products'
import { Container, Grid } from '@mantine/core'
import range from 'lodash/range'
import { useCallback } from 'react'

export default function Home() {
  const { data, isLoading } = useGetProducts(1, 10)
  const getProductLink = useCallback(
    (product: Product) => `/products/${product.id}`,
    []
  )

  return (
    <MainLayout>
      <Container size="lg">
        <Grid columns={4} gutter="xl">
          {isLoading
            ? range(10).map((index) => (
                <Grid.Col key={index} span={1}>
                  <ProductCardSkeleton />
                </Grid.Col>
              ))
            : data?.data.map((product) => (
                <Grid.Col key={product.id} span={1}>
                  <ProductCard
                    product={product}
                    getProductLink={getProductLink}
                  />
                </Grid.Col>
              ))}
        </Grid>
      </Container>
    </MainLayout>
  )
}
