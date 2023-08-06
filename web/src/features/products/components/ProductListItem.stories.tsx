import { Meta, StoryObj } from '@storybook/react'
import Decimal from 'decimal.js'
import { ProductListItem } from './ProductListItem'

const IMAGE_SIZE = 200

const meta: Meta<typeof ProductListItem> = {
  title: 'Features/Products/ProductListItem',
  component: ProductListItem,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof ProductListItem>

export const Default: Story = {
  args: {
    product: {
      id: '1',
      name: 'Product',
      description: 'Description',
      price: new Decimal(10.0),
      stockQuantity: 10,
      category: 'Category',
      seller: 'Seller',
      imageUrl: 'https://via.placeholder.com/200',
    },
    imageSize: IMAGE_SIZE,
  },
}

export const NotEditable: Story = {
  args: {
    product: {
      id: '1',
      name: 'Product',
      description: 'Description',
      price: new Decimal(10.0),
      stockQuantity: 10,
      category: 'Category',
      seller: 'Seller',
      imageUrl: 'https://via.placeholder.com/200',
    },
    imageSize: IMAGE_SIZE,
  },
}
