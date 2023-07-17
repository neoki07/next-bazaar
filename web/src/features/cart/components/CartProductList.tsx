import { Text, createStyles } from '@mantine/core'
import { range } from 'lodash'
import { CartProduct } from '../types'
import { CartProductInfo } from './CartProductInfo'
import { CartProductInfoSkeleton } from './CartProductInfoSkeleton'

const IMAGE_SIZE = 192

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
  isLoading?: boolean
}

export function CartProductList({
  cartProducts,
  isLoading,
}: CartProductListProps) {
  const { classes } = useStyles()

  return (
    <>
      {isLoading || cartProducts === undefined ? (
        <ul className={classes.list}>
          {range(3).map((index) => (
            <CartProductInfoSkeleton
              key={index}
              className={classes.listItem}
              imageSize={IMAGE_SIZE}
            />
          ))}
        </ul>
      ) : cartProducts.length > 0 ? (
        <ul className={classes.list}>
          {cartProducts.map((product) => (
            <CartProductInfo
              key={product.id}
              className={classes.listItem}
              cartProduct={product}
              imageSize={IMAGE_SIZE}
            />
          ))}
        </ul>
      ) : (
        <Text>No products</Text>
      )}
    </>
  )
}
