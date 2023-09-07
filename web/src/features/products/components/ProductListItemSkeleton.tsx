import { FixedSizeImage } from '@/components/Image'
import { PriceSkeleton } from '@/components/Price'
import { Flex, Skeleton, Stack, Text, clsx, createStyles } from '@mantine/core'

const useStyles = createStyles(() => ({
  root: {
    listStyle: 'none',
  },
}))

interface ProductListItemSkeletonProps {
  className?: string
  imageSize: number
}

export function ProductListItemSkeleton({
  className,
  imageSize,
}: ProductListItemSkeletonProps) {
  const { classes } = useStyles()

  return (
    <li className={clsx(classes.root, className)}>
      <Stack spacing={4}>
        <Flex gap="xs">
          <FixedSizeImage width={imageSize} height={imageSize} isLoading />
          <div style={{ width: '200px' }}>
            <Stack spacing="xs">
              <Skeleton>
                <Text fz="md">dummy</Text>
              </Skeleton>
              <PriceSkeleton width="50%" />
            </Stack>
          </div>
          <div />
        </Flex>
      </Stack>
    </li>
  )
}
