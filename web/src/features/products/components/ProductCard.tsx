import { ResponsiveSquareImage } from '@/components/Image'
import { Price } from '@/components/Price'
import { Product } from '@/features/products'
import { Group, Stack, Text, rem, useMantineTheme } from '@mantine/core'
import { IconBuildingStore } from '@tabler/icons-react'
import Link from 'next/link'

interface ProductCardProps {
  product: Product
  getProductLink: (product: Product) => string
}

export function ProductCard({ product, getProductLink }: ProductCardProps) {
  const theme = useMantineTheme()
  const productLink = getProductLink(product)

  return (
    <Stack spacing="xs">
      {product.imageUrl !== undefined && (
        <Link href={productLink}>
          <ResponsiveSquareImage src={product.imageUrl} alt={product.name} />
        </Link>
      )}
      <Stack spacing={4}>
        <Text
          size="xs"
          sx={(theme) => ({
            color: theme.colors.gray[7],
          })}
        >
          {product.category}
        </Text>
        <Link href={productLink}>
          <Text>{product.name}</Text>
        </Link>
        <Price price={product.price} size="xl" weight="bold" />
        <Group spacing={rem(3)} mt={rem(1)}>
          <IconBuildingStore
            size={16}
            color={theme.colors.gray[7]}
            strokeWidth={1.25}
          />
          <Text
            size="sm"
            sx={(theme) => ({
              color: theme.colors.gray[7],
            })}
          >
            {product.seller}
          </Text>
        </Group>
      </Stack>
    </Stack>
  )
}
