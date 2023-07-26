import { Meta, StoryObj } from '@storybook/react'
import { CartProductListItemSkeleton } from './CartProductListItemSkeleton'

const IMAGE_SIZE = 200

const meta: Meta<typeof CartProductListItemSkeleton> = {
  title: 'Features/Cart/CartProductListItemSkeleton',
  component: CartProductListItemSkeleton,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof CartProductListItemSkeleton>

export const Default: Story = {
  args: {
    imageSize: IMAGE_SIZE,
  },
}
