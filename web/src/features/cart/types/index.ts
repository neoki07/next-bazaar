import Decimal from 'decimal.js'

export interface CartProduct {
  id: string
  name: string
  description?: string
  price: Decimal
  quantity: number
  subtotal: Decimal
  imageUrl?: string
}

export interface Cart {
  products: CartProduct[]
  subtotal: Decimal
  shipping: Decimal
  tax: Decimal
  total: Decimal
}
