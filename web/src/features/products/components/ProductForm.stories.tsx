import { Meta, StoryObj } from '@storybook/react'
import { ProductForm } from './ProductForm'

const meta: Meta<typeof ProductForm> = {
  title: 'Features/Products/ProductForm',
  component: ProductForm,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof ProductForm>

export const Default: Story = {
  args: {
    allCategories: [
      { id: '1', name: 'T-Shirt' },
      { id: '2', name: 'Jeans' },
      { id: '3', name: 'TV' },
      { id: '4', name: 'Sofa' },
    ],
    onSubmit: (data) => alert(JSON.stringify(data, null, 2)),
  },
}
