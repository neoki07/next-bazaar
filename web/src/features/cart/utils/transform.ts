import {
  CartDomainCartProductResponse,
  CartDomainCartResponse,
} from '@/api/model'
import Decimal from 'decimal.js'
import { Cart, CartProduct } from '../types'

function transformCartProduct(
  cartProduct: CartDomainCartProductResponse
): CartProduct {
  if (cartProduct.id === undefined) {
    throw new Error(
      'required field `id` is undefined:' + JSON.stringify(cartProduct)
    )
  } else if (cartProduct.name === undefined) {
    throw new Error(
      'required field `name` is undefined:' + JSON.stringify(cartProduct)
    )
  } else if (cartProduct.price === undefined) {
    throw new Error(
      'required field `price` is undefined:' + JSON.stringify(cartProduct)
    )
  } else if (cartProduct.quantity === undefined) {
    throw new Error(
      'required field `quantity` is undefined:' + JSON.stringify(cartProduct)
    )
  } else if (cartProduct.subtotal === undefined) {
    throw new Error(
      'required field `subtotal` is undefined:' + JSON.stringify(cartProduct)
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
    throw new Error(
      'required field `products` is undefined:' + JSON.stringify(cart)
    )
  } else if (cart.subtotal === undefined) {
    throw new Error(
      'required field `subtotal` is undefined:' + JSON.stringify(cart)
    )
  } else if (cart.shipping === undefined) {
    throw new Error(
      'required field `shipping` is undefined:' + JSON.stringify(cart)
    )
  } else if (cart.tax === undefined) {
    throw new Error('required field `tax` is undefined:' + JSON.stringify(cart))
  } else if (cart.total === undefined) {
    throw new Error(
      'required field `total` is undefined:' + JSON.stringify(cart)
    )
  }

  return {
    products: cart.products.map((product) => transformCartProduct(product)),
    subtotal: new Decimal(cart.subtotal),
    shipping: new Decimal(cart.shipping),
    tax: new Decimal(cart.tax),
    total: new Decimal(cart.total),
  }
}
