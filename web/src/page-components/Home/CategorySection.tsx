import {
  NOTIFY_NOT_IMPLEMENTED_ERRORS,
  notifyNotImplementedError,
} from '@/features/notification/not-implemented'
import {
  Category,
  Product,
  ProductCard,
  ProductCardSkeleton,
  useGetProducts,
} from '@/features/products'
import { Button, Center, SimpleGrid, Stack, Title } from '@mantine/core'
import { range } from 'lodash'

interface CategorySectionProps {
  category: Category
  getProductLink: (product: Product) => string
  productCount: number
}

export function CategorySection({
  category,
  getProductLink,
  productCount,
}: CategorySectionProps) {
  const { data, isLoading } = useGetProducts({
    page: 1,
    pageSize: productCount,
    categoryId: category.id,
  })

  return (
    <Stack spacing="xl">
      <Title order={2}>{category.name}</Title>
      <SimpleGrid cols={4} spacing="xl">
        {isLoading
          ? range(productCount).map((index) => (
              <ProductCardSkeleton key={index} />
            ))
          : data?.data.map((product) => (
              <ProductCard
                key={product.id}
                product={product}
                getProductLink={getProductLink}
              />
            ))}
      </SimpleGrid>
      <Center>
        {/* TODO: Add link to view more products */}
        {/* <Link href="#"> */}
        <Button
          variant="default"
          fullWidth
          onClick={() =>
            notifyNotImplementedError(
              NOTIFY_NOT_IMPLEMENTED_ERRORS.ViewMoreProducts
            )
          }
        >
          View More
        </Button>
        {/* </Link> */}
      </Center>
    </Stack>
  )
}
