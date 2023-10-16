import { zodResolver } from '@hookform/resolvers/zod'
import type { Meta, StoryObj } from '@storybook/react'
import { z } from 'zod'
import { renderDecorator } from '../utils/storybook'
import { RadioGroup } from './RadioGroup'

const label = 'Browser'
const name = 'browser'
const options = [
  { label: 'Firefox', value: 'firefox' },
  { label: 'Edge', value: 'edge' },
  { label: 'Chrome', value: 'chrome' },
  { label: 'Opera', value: 'opera' },
  { label: 'Safari', value: 'safari' },
]

const schema = z.object({
  browser: z.string().array().min(1, { message: 'Required' }),
})

const resolver = zodResolver(schema)

const defaultValues = {
  browser: [],
}

const meta: Meta<typeof RadioGroup> = {
  title: 'Example/Form/RadioGroup',
  component: RadioGroup,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof RadioGroup>

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
