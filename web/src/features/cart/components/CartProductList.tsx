import { Text, createStyles } from '@mantine/core'
import { range } from 'lodash'
import { CartProduct } from '../types'
import { CartProductListItem } from './CartProductListItem'
import { CartProductListItemSkeleton } from './CartProductListItemSkeleton'

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

interface CartProductListProps {
  cartProducts?: CartProduct[]
  imageSize: number
  isLoading?: boolean
  itemsCountOnLoad?: number
  editable?: boolean
  onChangeQuantity?: (id: string, quantity: number) => void
  onDelete?: (id: string) => void
}

export function CartProductList({
  cartProducts,
  imageSize,
  isLoading,
  itemsCountOnLoad = 3,
  editable = true,
  onChangeQuantity,
  onDelete,
}: CartProductListProps) {
  const { classes } = useStyles()

  return (
    <>
      {isLoading || cartProducts === undefined ? (
        <ul className={classes.list}>
          {range(itemsCountOnLoad).map((index) => (
            <CartProductListItemSkeleton
              key={index}
              className={classes.listItem}
              imageSize={imageSize}
            />
          ))}
        </ul>
      ) : cartProducts.length > 0 ? (
        <ul className={classes.list}>
          {cartProducts.map((product) => (
            <CartProductListItem
              key={product.id}
              className={classes.listItem}
              cartProduct={product}
              imageSize={imageSize}
              editable={editable}
              onChangeQuantity={onChangeQuantity}
              onDelete={onDelete}
            />
          ))}
        </ul>
      ) : (
        <Text>No products</Text>
      )}
    </>
  )
}
