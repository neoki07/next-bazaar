import { DashboardLayout } from '@/features/dashboard'
import { ProductList, useGetProducts } from '@/features/products'
import { Button, Center, Pagination, Stack, Title } from '@mantine/core'
import { IconPlus } from '@tabler/icons-react'
import { useState } from 'react'

export function Products() {
  const [page, setPage] = useState(1)

  const { data: products } = useGetProducts({
    page,
    pageSize: 10,
  })

  return (
    <DashboardLayout>
      <Stack spacing="xl">
        <Title order={1}>Your Products</Title>
        <div>
          <Button leftIcon={<IconPlus size="1rem" />}>
            Register New Product
          </Button>
        </div>
        <ProductList products={products?.data} imageSize={120} />
        <Center>
          <Pagination
            value={page}
            onChange={setPage}
            total={products?.meta.page_count ?? 0}
          />
        </Center>
      </Stack>
    </DashboardLayout>
  )
}
