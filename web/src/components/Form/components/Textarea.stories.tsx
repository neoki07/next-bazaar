import { zodResolver } from '@hookform/resolvers/zod'
import type { Meta, StoryObj } from '@storybook/react'
import { z } from 'zod'
import { renderDecorator } from '../utils/storybook'
import { Textarea } from './Textarea'

const label = 'Comment'
const name = 'comment'

const schema = z.object({
  comment: z.string().min(1, { message: 'Required' }),
})

const resolver = zodResolver(schema)

const defaultValues = {
  comment: '',
}

const meta: Meta<typeof Textarea> = {
  title: 'Example/Form/Textarea',
  component: Textarea,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof Textarea>

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
