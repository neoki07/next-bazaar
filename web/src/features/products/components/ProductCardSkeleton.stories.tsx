import { Meta, StoryObj } from '@storybook/react'
import { ProductCardSkeleton } from './ProductCardSkeleton'

const meta: Meta<typeof ProductCardSkeleton> = {
  title: 'Features/Products/ProductCardSkeleton',
  component: ProductCardSkeleton,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof ProductCardSkeleton>

export const Default: Story = {
  args: {
    imageSize: 300,
  },
}
