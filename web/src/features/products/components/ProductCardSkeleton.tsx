import { ResponsiveSquareImage } from '@/components/Image'
import { PriceSkeleton } from '@/components/Price'
import { Skeleton, Stack, Text } from '@mantine/core'

export function ProductCardSkeleton() {
  return (
    <Stack spacing="xs">
      <ResponsiveSquareImage isLoading />
      <Stack spacing={4}>
        <Skeleton visible width="50%">
          <Text size="xs">dummy</Text>
        </Skeleton>
        <Skeleton visible>
          <Text>dummy</Text>
        </Skeleton>
        <PriceSkeleton width="50%" />
        <Skeleton width="50%">
          <Text size="sm">dummy</Text>
        </Skeleton>
      </Stack>
    </Stack>
  )
}
