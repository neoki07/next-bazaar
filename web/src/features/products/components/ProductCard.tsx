import { Price } from '@/components/Price'
import { Product } from '@/features/products'
import { Group, Image, Stack, Text, rem, useMantineTheme } from '@mantine/core'
import { IconBuildingStore } from '@tabler/icons-react'

type ProductCardProps = {
  product: Product
}

export function ProductCard({ product }: ProductCardProps) {
  const theme = useMantineTheme()

  return (
    <Stack spacing="xs">
      <Image src={product.imageUrl} alt="Tesla Model S" />
      <div>
        <Text
          size="xs"
          sx={(theme) => ({
            color: theme.colors.gray[7],
          })}
        >
          {product.category}
        </Text>
        <Text>{product.name}</Text>
        <Price price={product.price} />
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
      </div>
    </Stack>
  )
}
