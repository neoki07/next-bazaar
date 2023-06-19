import { Skeleton, Text } from '@mantine/core'

interface PriceSkeletonProps {
  width?: string | number | undefined
}

export function PriceSkeleton({ width }: PriceSkeletonProps) {
  return (
    <Skeleton visible width={width}>
      <Text size="xl" weight="bold">
        dummy
      </Text>
    </Skeleton>
  )
}
