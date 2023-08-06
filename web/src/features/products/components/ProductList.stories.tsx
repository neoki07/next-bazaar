import { Meta, StoryObj } from '@storybook/react'
import Decimal from 'decimal.js'
import { ProductList } from './ProductList'

const IMAGE_SIZE = 200

const meta: Meta<typeof ProductList> = {
  title: 'Features/Products/ProductList',
  component: ProductList,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof ProductList>

export const Default: Story = {
  args: {
    products: [
      {
        id: '1',
        name: 'Product 1',
        description: 'Description 1',
        price: new Decimal(10.0),
        stockQuantity: 10,
        category: 'Category 1',
        seller: 'Seller 1',
        imageUrl: 'https://via.placeholder.com/200',
      },
      {
        id: '2',
        name: 'Product 2',
        description: 'Description 2',
        price: new Decimal(20.0),
        stockQuantity: 10,
        category: 'Category 2',
        seller: 'Seller 2',
        imageUrl: 'https://via.placeholder.com/200',
      },
    ],
    imageSize: IMAGE_SIZE,
  },
}

export const NotEditable: Story = {
  args: {
    products: [
      {
        id: '1',
        name: 'Product 1',
        description: 'Description 1',
        price: new Decimal(10.0),
        stockQuantity: 10,
        category: 'Category 1',
        seller: 'Seller 1',
        imageUrl: 'https://via.placeholder.com/200',
      },
      {
        id: '2',
        name: 'Product 2',
        description: 'Description 2',
        price: new Decimal(20.0),
        stockQuantity: 10,
        category: 'Category 2',
        seller: 'Seller 2',
        imageUrl: 'https://via.placeholder.com/200',
      },
    ],
    imageSize: IMAGE_SIZE,
  },
}

export const Loading: Story = {
  args: {
    imageSize: IMAGE_SIZE,
    isLoading: true,
  },
}
