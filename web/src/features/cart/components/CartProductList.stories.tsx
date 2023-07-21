import { Meta, StoryObj } from '@storybook/react'
import Decimal from 'decimal.js'
import { CartProductList } from './CartProductList'

const IMAGE_SIZE = 200

const meta: Meta<typeof CartProductList> = {
  title: 'Features/Cart/CartProductList',
  component: CartProductList,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof CartProductList>

export const Default: Story = {
  args: {
    cartProducts: [
      {
        id: '1',
        name: 'Product 1',
        description: 'Description 1',
        price: new Decimal(10.0),
        quantity: 1,
        subtotal: new Decimal(10.0),
        imageUrl: 'https://via.placeholder.com/200',
      },
      {
        id: '2',
        name: 'Product 2',
        description: 'Description 2',
        price: new Decimal(20.0),
        quantity: 2,
        subtotal: new Decimal(40.0),
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
