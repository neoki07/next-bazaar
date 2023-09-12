import { FixedSizeImage } from '@/components/Image'
import { PriceSkeleton } from '@/components/Price'
import {
  CloseButton,
  Flex,
  Skeleton,
  Stack,
  Text,
  clsx,
  createStyles,
  rem,
} from '@mantine/core'

const useStyles = createStyles(() => ({
  root: {
    listStyle: 'none',
  },
}))

interface CartProductListItemSkeletonProps {
  className?: string
  imageSize: number
}

export function CartProductListItemSkeleton({
  className,
  imageSize,
}: CartProductListItemSkeletonProps) {
  const { classes } = useStyles()

  return (
    <li className={clsx(classes.root, className)}>
      <Stack spacing={4}>
        <Flex gap="xs">
          <FixedSizeImage width={imageSize} height={imageSize} isLoading />
          <Stack spacing="xs" style={{ flex: 1 }}>
            <Skeleton width={rem(200)}>
              <Text fz="md">dummy</Text>
            </Skeleton>
            <PriceSkeleton width={rem(100)} />
            <Skeleton width={rem(100)}>
              <Text fz="md">dummy</Text>
            </Skeleton>
          </Stack>
          <div>
            <Skeleton>
              <CloseButton />
            </Skeleton>
          </div>
        </Flex>
      </Stack>
    </li>
  )
}
