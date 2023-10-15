import { zodResolver } from '@hookform/resolvers/zod'
import type { Meta, StoryObj } from '@storybook/react'
import { z } from 'zod'
import { renderDecorator } from '../utils/storybook'
import { NativeNumberSelect } from './NativeNumberSelect'

const label = 'Amount'
const name = 'amount'
const options = [1, 2, 3, 4, 5]

const schema = z.object({
  amount: z.number().min(2),
})

const resolver = zodResolver(schema)

const defaultValues = {
  amount: 1,
}

const meta: Meta<typeof NativeNumberSelect> = {
  title: 'Example/Form/NativeNumberSelect',
  component: NativeNumberSelect,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof NativeNumberSelect>

export const Default: Story = {
  args: {
    label,
    name,
    options,
  },
  decorators: [(Story) => renderDecorator(Story, resolver, defaultValues)],
}

export const WithAsterisk: Story = {
  args: {
    label,
    name,
    options,
    withAsterisk: true,
  },
  decorators: [(Story) => renderDecorator(Story, resolver, defaultValues)],
}

export const ErrorMessage: Story = {
  args: {
    label,
    name,
    options,
    withAsterisk: true,
  },
  decorators: [
    (Story) => renderDecorator(Story, resolver, defaultValues, true),
  ],
}
