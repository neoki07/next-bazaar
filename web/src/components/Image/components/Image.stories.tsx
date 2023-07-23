import { Meta, StoryObj } from '@storybook/react'
import { Image } from './Image'

const IMAGE_SIZE = 200

const meta: Meta<typeof Image> = {
  title: 'Components/Image/Image',
  component: Image,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof Image>

export const Default: Story = {
  args: {
    src: 'https://via.placeholder.com/200',
    alt: 'Example Image',
    width: IMAGE_SIZE,
    height: IMAGE_SIZE,
  },
}

export const Loading: Story = {
  args: {
    src: 'https://via.placeholder.com/200',
    alt: 'Example Image',
    width: IMAGE_SIZE,
    height: IMAGE_SIZE,
    isLoading: true,
  },
}
