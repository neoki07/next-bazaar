import { Meta, StoryObj } from '@storybook/react'
import Decimal from 'decimal.js'
import { ProductCard } from './ProductCard'

const meta: Meta<typeof ProductCard> = {
  title: 'Features/Products/ProductCard',
  component: ProductCard,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof ProductCard>

export const Default: Story = {
  args: {
    product: {
      id: '1',
      name: 'Product',
      category: 'Category',
      price: new Decimal(10.0),
      stockQuantity: 5,
      seller: 'Seller',
      imageUrl: 'https://via.placeholder.com/300',
    },
    getProductLink: () => '#',
    imageSize: 300,
  },
}
