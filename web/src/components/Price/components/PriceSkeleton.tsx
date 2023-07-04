import { Skeleton, Text } from '@mantine/core'

interface PriceSkeletonProps {
  width?: string | number | undefined
  className?: string
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | string | number
  weight?: React.CSSProperties['fontWeight']
}

export function PriceSkeleton({
  width,
  className,
  size,
  weight,
}: PriceSkeletonProps) {
  return (
    <Skeleton visible width={width}>
      <Text className={className} size={size} weight={weight}>
        dummy
      </Text>
    </Skeleton>
  )
}
