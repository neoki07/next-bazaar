import { zodResolver } from '@hookform/resolvers/zod'
import type { Meta, StoryObj } from '@storybook/react'
import { z } from 'zod'
import { renderDecorator } from '../utils/storybook'
import { PasswordInput } from './PasswordInput'

const label = 'Password'
const name = 'password'

const schema = z.object({
  password: z.string().min(1, { message: 'Required' }),
})

const resolver = zodResolver(schema)

const defaultValues = {
  password: '',
}

const meta: Meta<typeof PasswordInput> = {
  title: 'Components/Form/PasswordInput',
  component: PasswordInput,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof PasswordInput>

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
