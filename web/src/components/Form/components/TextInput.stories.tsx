import { zodResolver } from '@hookform/resolvers/zod'
import type { Meta, StoryObj } from '@storybook/react'
import { z } from 'zod'
import { renderDecorator } from '../utils/storybook'
import { TextInput } from './TextInput'

const label = 'Username'
const name = 'username'

const schema = z.object({
  username: z.string().min(1, { message: 'Required' }),
})

const resolver = zodResolver(schema)

const defaultValues = {
  username: '',
}

const meta: Meta<typeof TextInput> = {
  title: 'Example/Form/TextInput',
  component: TextInput,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof TextInput>

export const Default: Story = {
  args: {
    label,
    name,
  },
  decorators: [(Story) => renderDecorator(Story, resolver, defaultValues)],
}

export const WithAsterisk: Story = {
  args: {
    label,
    name,
    withAsterisk: true,
  },
  decorators: [(Story) => renderDecorator(Story, resolver, defaultValues)],
}

export const ErrorMessage: Story = {
  args: {
    label,
    name,
    withAsterisk: true,
  },
  decorators: [
    (Story) => renderDecorator(Story, resolver, defaultValues, true),
  ],
}
