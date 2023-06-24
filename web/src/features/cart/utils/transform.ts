import {
  CartDomainCartProductResponse,
  CartDomainCartResponse,
} from '@/api/model'
import Decimal from 'decimal.js'
import { Cart, CartProduct } from '../types'

function transformCartProduct(
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

export function transformCart(cart: CartDomainCartResponse): Cart {
  if (cart.products === undefined) {
    throw new Error('required fields are undefined:' + JSON.stringify(cart))
  }

  return {
    products: cart.products.map((product) => transformCartProduct(product)),
  }
}
