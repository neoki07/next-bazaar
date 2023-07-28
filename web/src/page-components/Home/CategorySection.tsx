import {
  Category,
  Product,
  ProductCard,
  ProductCardSkeleton,
  useGetProducts,
} from '@/features/products'
import {
  Button,
  Center,
  Grid,
  Stack,
  Title,
  createStyles,
  rem,
} from '@mantine/core'
import range from 'lodash/range'
import Link from 'next/link'

const PAGE_SIZE = 8

const useStyles = createStyles((theme) => ({
  viewMoreButton: {
    width: rem(400),
    height: rem(40),
    color: theme.colors.gray[7],
  },
}))

interface CategorySectionProps {
  category: Category
  getProductLink: (product: Product) => string
  imageSize: number
}

export function CategorySection({
  category,
  getProductLink,
  imageSize,
}: CategorySectionProps) {
  const { classes } = useStyles()
  const { data, isLoading } = useGetProducts(1, PAGE_SIZE, category.id)

  return (
    <Stack spacing="xl">
      <Title order={2}>{category.name}</Title>
      <Grid columns={4} gutter="xl">
        {isLoading
          ? range(PAGE_SIZE).map((index) => (
              <Grid.Col key={index} span={1}>
                <ProductCardSkeleton imageSize={imageSize} />
              </Grid.Col>
            ))
          : data?.data.map((product) => (
              <Grid.Col key={product.id} span={1}>
                <ProductCard
                  product={product}
                  getProductLink={getProductLink}
                  imageSize={imageSize}
                />
              </Grid.Col>
            ))}
      </Grid>
      <Center>
        <Link href="#">
          <Button className={classes.viewMoreButton} variant="default">
            View More
          </Button>
        </Link>
      </Center>
    </Stack>
  )
}
