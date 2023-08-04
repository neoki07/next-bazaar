import { ProductCardSkeleton } from '@/features/products'
import {
  Button,
  Center,
  Container,
  Grid,
  Skeleton,
  Stack,
  Title,
  createStyles,
  rem,
} from '@mantine/core'
import range from 'lodash/range'

const useStyles = createStyles((theme) => ({
  viewMoreButton: {
    width: rem(400),
    height: rem(40),
  },
}))

interface CategorySectionSkeletonProps {
  imageSize: number
  productCount: number
}

export function CategorySectionSkeleton({
  imageSize,
  productCount,
}: CategorySectionSkeletonProps) {
  const { classes } = useStyles()

  return (
    <Container size="lg">
      <Stack spacing="xl">
        <Skeleton width="15%">
          <Title order={2}>dummy</Title>
        </Skeleton>
        <Grid columns={4} gutter="xl">
          {range(productCount).map((index) => (
            <Grid.Col key={index} span={1}>
              <ProductCardSkeleton imageSize={imageSize} />
            </Grid.Col>
          ))}
        </Grid>
        <Center>
          <Skeleton>
            <Button className={classes.viewMoreButton} variant="default">
              View More
            </Button>
          </Skeleton>
        </Center>
      </Stack>
    </Container>
  )
}
