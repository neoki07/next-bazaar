import { Meta, StoryObj } from '@storybook/react'
import Decimal from 'decimal.js'
import { CartProductListItem } from './CartProductListItem'

const IMAGE_SIZE = 200

const meta: Meta<typeof CartProductListItem> = {
  title: 'Features/Cart/CartProductListItem',
  component: CartProductListItem,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof CartProductListItem>

export const Default: Story = {
  args: {
    cartProduct: {
      id: '1',
      name: 'Product',
      description: 'Description',
      price: new Decimal(10.0),
      quantity: 5,
      subtotal: new Decimal(50.0),
      imageUrl: 'https://via.placeholder.com/200',
    },
    imageSize: IMAGE_SIZE,
  },
}

export const NotEditable: Story = {
  args: {
    cartProduct: {
      id: '1',
      name: 'Product',
      description: 'Description',
      price: new Decimal(10.0),
      quantity: 5,
      subtotal: new Decimal(50.0),
      imageUrl: 'https://via.placeholder.com/200',
    },
    imageSize: IMAGE_SIZE,
    editable: false,
  },
}
