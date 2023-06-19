import { Image } from '@/components/Image'
import { PriceSkeleton } from '@/components/Price'
import { Flex, Group, Skeleton, Stack, Text } from '@mantine/core'

interface CartProductInfoSkeletonProps {
  imageSize: number
}

export function CartProductInfoSkeleton({
  imageSize,
}: CartProductInfoSkeletonProps) {
  return (
    <Group my="sm">
      <Stack spacing={4}>
        <Flex gap="xs">
          <Image isLoading alt="" width={imageSize} height={imageSize} />
          <div style={{ width: '200px' }}>
            <Stack spacing="xs">
              <Skeleton>
                <Text fz="md">dummy</Text>
              </Skeleton>
              <PriceSkeleton width="50%" />
              <Skeleton width="50%" height={60} />
            </Stack>
          </div>
          <div />
        </Flex>
      </Stack>
    </Group>
  )
}
