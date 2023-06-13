import { Divider, Text } from '@mantine/core'
import { Fragment } from 'react'
import { CartProduct } from '../types'
import { CartProductInfo } from './CartProductInfo'

interface CartProductListProps {
  cartProducts: CartProduct[]
}

export function CartProductList({ cartProducts }: CartProductListProps) {
  return (
    <>
      {cartProducts.length > 0 ? (
        <>
          <Divider mb="xl" />
          {cartProducts.map((product, index) => (
            <Fragment key={product.id}>
              {index !== 0 && <Divider my="xl" />}
              <CartProductInfo cartProduct={product} />
            </Fragment>
          ))}
          <Divider mt="xl" />
        </>
      ) : (
        <Text>No products</Text>
      )}
    </>
  )
}
