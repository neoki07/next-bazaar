import { Divider, Text } from '@mantine/core'
import { range } from 'lodash'
import { Fragment } from 'react'
import { CartProduct } from '../types'
import { CartProductInfo } from './CartProductInfo'
import { CartProductInfoSkeleton } from './CartProductInfoSkeleton'

const IMAGE_SIZE = 192

interface CartProductListProps {
  cartProducts?: CartProduct[]
  isLoading?: boolean
}

export function CartProductList({
  cartProducts,
  isLoading,
}: CartProductListProps) {
  return (
    <>
      {isLoading || cartProducts === undefined ? (
        <>
          <Divider mb="xl" />
          {range(3).map((index) => (
            <Fragment key={index}>
              {index !== 0 && <Divider my="xl" />}
              <CartProductInfoSkeleton imageSize={IMAGE_SIZE} />
            </Fragment>
          ))}
          <Divider mt="xl" />
        </>
      ) : cartProducts.length > 0 ? (
        <>
          <Divider mb="xl" />
          {cartProducts.map((product, index) => (
            <Fragment key={product.id}>
              {index !== 0 && <Divider my="xl" />}
              <CartProductInfo cartProduct={product} imageSize={IMAGE_SIZE} />
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
