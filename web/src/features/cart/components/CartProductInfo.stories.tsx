import { Meta, StoryObj } from '@storybook/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import Decimal from 'decimal.js'
import { CartProductInfo } from './CartProductInfo'

const queryClient = new QueryClient()

const meta: Meta<typeof CartProductInfo> = {
  title: 'Features/Cart/CartProductInfo',
  component: CartProductInfo,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof CartProductInfo>

export const Default: Story = {
  args: {
    cartProduct: {
      id: '1',
      name: 'Product',
      description: 'Description',
      price: new Decimal(10.0),
      quantity: 5,
      subtotal: new Decimal(50.0),
      imageUrl: 'https://via.placeholder.com/300',
    },
    imageSize: 300,
  },
  decorators: [
    (Story) => (
      <QueryClientProvider client={queryClient}>
        <Story />
      </QueryClientProvider>
    ),
  ],
}
