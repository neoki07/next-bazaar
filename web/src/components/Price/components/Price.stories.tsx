import { Meta, StoryObj } from '@storybook/react'
import Decimal from 'decimal.js'
import { Price } from './Price'

const meta: Meta<typeof Price> = {
  title: 'Components/Price/Price',
  component: Price,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof Price>

export const Default: Story = {
  args: {
    price: new Decimal(12.34),
  },
}

export const XLAndBold: Story = {
  args: {
    price: new Decimal(12.34),
    size: 'xl',
    weight: 'bold',
  },
}
