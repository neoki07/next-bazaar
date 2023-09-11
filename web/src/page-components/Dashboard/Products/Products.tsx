import { DashboardLayout } from '@/features/dashboard'
import { ProductList, useGetProducts } from '@/features/products'
import { Button, Center, Pagination, Stack, Title } from '@mantine/core'
import { IconPlus } from '@tabler/icons-react'
import { useRouter } from 'next/router'
import { useCallback, useState } from 'react'

export function Products() {
  const router = useRouter()

  const [page, setPage] = useState(1)

  const { data: products } = useGetProducts({
    page,
    pageSize: 10,
  })

  const goToAddNewProductPage = useCallback(() => {
    router.push('/dashboard/products/new')
  }, [router])

  return (
    <DashboardLayout>
      <Stack spacing="xl">
        <Title order={1}>Your Products</Title>
        <div>
          <Button
            color="dark"
            leftIcon={<IconPlus size="1rem" />}
            onClick={goToAddNewProductPage}
          >
            Add New Product
          </Button>
        </div>
        <ProductList products={products?.data} imageSize={120} />
        <Center>
          <Pagination
            color="dark"
            value={page}
            onChange={setPage}
            total={products?.meta.page_count ?? 0}
          />
        </Center>
      </Stack>
    </DashboardLayout>
  )
}
