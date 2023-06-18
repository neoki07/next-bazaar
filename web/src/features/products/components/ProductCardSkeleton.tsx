import { PriceSkeleton } from '@/components/Price'
import { Skeleton, Stack, Text, useMantineTheme } from '@mantine/core'

export function ProductCardSkeleton() {
  const theme = useMantineTheme()

  return (
    <Stack spacing="xs">
      <Skeleton visible>
        {/* TODO: improve the way height is decided */}
        <svg viewBox="0 0 260 253.19" />
      </Skeleton>
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
