import { Meta, StoryObj } from '@storybook/react'
import { ProductListItemSkeleton } from './ProductListItemSkeleton'

const IMAGE_SIZE = 200

const meta: Meta<typeof ProductListItemSkeleton> = {
  title: 'Features/Products/ProductListItemSkeleton',
  component: ProductListItemSkeleton,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof ProductListItemSkeleton>

export const Default: Story = {
  args: {
    imageSize: IMAGE_SIZE,
  },
}
