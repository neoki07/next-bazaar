import { Meta, StoryObj } from '@storybook/react'
import { FixedSizeImage } from './FixedSizeImage'

const src = 'https://via.placeholder.com/300x200'
const alt = 'Example Image'
const width = 300
const height = 200

const meta: Meta<typeof FixedSizeImage> = {
  title: 'Components/Image/FixedSizeImage',
  component: FixedSizeImage,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof FixedSizeImage>

export const Default: Story = {
  args: {
    src,
    alt,
    width,
    height,
  },
}

export const Loading: Story = {
  args: {
    src,
    alt,
    width,
    height,
    isLoading: true,
  },
}
