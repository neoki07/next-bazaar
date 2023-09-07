import { FixedSizeImage } from '@/components/Image'
import { Price } from '@/components/Price'
import {
  Flex,
  Stack,
  Text,
  UnstyledButton,
  clsx,
  createStyles,
} from '@mantine/core'
import { IconPencil } from '@tabler/icons-react'
import Link from 'next/link'
import { Product } from '../types'

const useStyles = createStyles((theme) => ({
  root: {
    listStyle: 'none',
  },
  editButton: {
    color: theme.colors.gray[6],

    '&:hover': {
      color: theme.colors.gray[8],
    },
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
            <FixedSizeImage
              src={product.imageUrl}
              alt={product.name}
              width={imageSize}
              height={imageSize}
            />
          )}
          <Stack spacing="xs" style={{ flex: 1 }}>
            <Text>{product.name}</Text>
            <Price price={product.price} size="xl" weight="bold" />
            <Text fz="sm">Stock Quantity: {product.stockQuantity}</Text>
          </Stack>
          <div>
            <Link href={`/dashboard/products/${product.id}/edit`}>
              <UnstyledButton className={classes.editButton}>
                <IconPencil size="1rem" />
              </UnstyledButton>
            </Link>
          </div>
        </Flex>
      </Stack>
    </li>
  )
}
