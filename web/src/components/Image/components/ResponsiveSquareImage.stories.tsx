import { Meta, StoryObj } from '@storybook/react'
import { ResponsiveSquareImage } from './ResponsiveSquareImage'

const src = 'https://via.placeholder.com/300x200'
const alt = 'Example Image'

const meta: Meta<typeof ResponsiveSquareImage> = {
  title: 'Components/Image/ResponsiveSquareImage',
  component: ResponsiveSquareImage,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof ResponsiveSquareImage>

export const Default: Story = {
  args: {
    src,
    alt,
  },
}

export const Loading: Story = {
  args: {
    src,
    alt,
    isLoading: true,
  },
}
