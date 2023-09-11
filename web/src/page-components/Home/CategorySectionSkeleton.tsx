import { ProductCardSkeleton } from '@/features/products'
import {
  Button,
  Center,
  SimpleGrid,
  Skeleton,
  Stack,
  Title,
} from '@mantine/core'
import range from 'lodash/range'

interface CategorySectionSkeletonProps {
  productCount: number
}

export function CategorySectionSkeleton({
  productCount,
}: CategorySectionSkeletonProps) {
  return (
    <Stack spacing="xl">
      <Skeleton width="15%">
        <Title order={2}>dummy</Title>
      </Skeleton>
      <SimpleGrid
        cols={4}
        spacing="xl"
        breakpoints={[
          { maxWidth: 'md', cols: 3 },
          { maxWidth: 'sm', cols: 2 },
          { maxWidth: 'xs', cols: 1 },
        ]}
      >
        {range(productCount).map((index) => (
          <ProductCardSkeleton key={index} />
        ))}
      </SimpleGrid>
      <Center>
        <Skeleton>
          <Button variant="default" fullWidth>
            View More
          </Button>
        </Skeleton>
      </Center>
    </Stack>
  )
}
