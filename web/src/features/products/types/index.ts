import Decimal from 'decimal.js'

export interface Product {
  id: string
  name: string
  description?: string
  price: Decimal
  stockQuantity: number
  categoryId: string
  category: string
  seller: string
  imageUrl?: string
}

export interface Category {
  id: string
  name: string
}
