import { CartDomainCartProductResponse } from '@/api/model'
import Decimal from 'decimal.js'
import { CartProduct } from '../types'

export function transformCartProduct(
  cartProduct: CartDomainCartProductResponse
): CartProduct {
  if (
    cartProduct.id === undefined ||
    cartProduct.name === undefined ||
    cartProduct.price === undefined ||
    cartProduct.quantity === undefined ||
    cartProduct.subtotal === undefined
  ) {
    throw new Error(
      'required fields are undefined:' + JSON.stringify(cartProduct)
    )
  }

  return {
    id: cartProduct.id,
    name: cartProduct.name,
    description: cartProduct.description,
    price: new Decimal(cartProduct.price),
    quantity: cartProduct.quantity,
    subtotal: new Decimal(cartProduct.subtotal),
    imageUrl: cartProduct.image_url,
  }
}
