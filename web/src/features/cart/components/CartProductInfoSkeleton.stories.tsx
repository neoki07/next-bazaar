import { Meta, StoryObj } from '@storybook/react'
import { CartProductInfoSkeleton } from './CartProductInfoSkeleton'

const IMAGE_SIZE = 200

const meta: Meta<typeof CartProductInfoSkeleton> = {
  title: 'Features/Cart/CartProductInfoSkeleton',
  component: CartProductInfoSkeleton,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof CartProductInfoSkeleton>

export const Default: Story = {
  args: {
    imageSize: IMAGE_SIZE,
  },
}
