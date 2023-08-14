import { Text, createStyles } from '@mantine/core'
import { range } from 'lodash'
import { Product } from '../types'
import { ProductListItem } from './ProductListItem'
import { ProductListItemSkeleton } from './ProductListItemSkeleton'

const useStyles = createStyles((theme) => ({
  list: {
    paddingLeft: 0,
    marginTop: 0,
    marginBottom: 0,
  },

  listItem: {
    paddingTop: theme.spacing.md,
    paddingBottom: theme.spacing.md,
    borderTop: `1px solid ${theme.colors.gray[3]}`,

    '&:last-child': {
      borderBottom: `1px solid ${theme.colors.gray[3]}`,
    },
  },
}))

interface ProductListProps {
  products?: Product[]
  imageSize: number
  isLoading?: boolean
  itemsCountOnLoad?: number
}

export function ProductList({
  products,
  imageSize,
  isLoading,
  itemsCountOnLoad = 3,
}: ProductListProps) {
  const { classes } = useStyles()

  return (
    <>
      {isLoading || products === undefined ? (
        <ul className={classes.list}>
          {range(itemsCountOnLoad).map((index) => (
            <ProductListItemSkeleton
              key={index}
              className={classes.listItem}
              imageSize={imageSize}
            />
          ))}
        </ul>
      ) : products.length > 0 ? (
        <ul className={classes.list}>
          {products.map((product) => (
            <ProductListItem
              key={product.id}
              className={classes.listItem}
              product={product}
              imageSize={imageSize}
            />
          ))}
        </ul>
      ) : (
        <Text>No products</Text>
      )}
    </>
  )
}
