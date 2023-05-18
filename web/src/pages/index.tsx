import { MainLayout } from '@/components/Layout'
import { Product, useGetProducts } from '@/features/products'
import { ProductCard } from '@/features/products/components/ProductCard'
import { Container, Grid } from '@mantine/core'
import { useCallback } from 'react'

export default function Home() {
  const { data, isLoading } = useGetProducts(1, 10)
  const getProductLink = useCallback(
    (product: Product) => `/products/${product.id}`,
    []
  )

  if (isLoading) {
    return <div>Loading...</div>
  }

  return (
    <MainLayout>
      <Container size="lg">
        <Grid columns={4} gutter="xl">
          {data?.data.map((product) => (
            <Grid.Col key={product.id} span={1}>
              <ProductCard product={product} getProductLink={getProductLink} />
            </Grid.Col>
          ))}
        </Grid>
      </Container>
    </MainLayout>
  )
}
