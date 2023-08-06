import { Image } from '@/components/Image'
import { Price } from '@/components/Price'
import { Flex, Stack, Text, clsx, createStyles } from '@mantine/core'
import { Product } from '../types'

const useStyles = createStyles(() => ({
  root: {
    listStyle: 'none',
  },
}))

interface ProductListItemProps {
  className?: string
  product: Product
  imageSize: number
}

export function ProductListItem({
  className,
  product,
  imageSize,
}: ProductListItemProps) {
  const { classes } = useStyles()

  return (
    <li className={clsx(classes.root, className)}>
      <Stack spacing={4}>
        <Flex gap="xs">
          {product.imageUrl !== undefined && (
            <Image
              src={product.imageUrl}
              alt={product.name}
              width={imageSize}
              height={imageSize}
            />
          )}
          <Stack spacing="xs">
            <Text fz="md">{product.name}</Text>
            <Price price={product.price} size="xl" weight="bold" />
          </Stack>
        </Flex>
      </Stack>
    </li>
  )
}
