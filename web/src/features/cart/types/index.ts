import Decimal from 'decimal.js'

export type Product = {
  id: string
  name: string
  price: number
  quantity: number
  isSelected: boolean
}

export type FormValues = {
  products: Product[]
}

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
}
