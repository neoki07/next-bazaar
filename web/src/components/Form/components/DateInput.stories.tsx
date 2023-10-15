import { zodResolver } from '@hookform/resolvers/zod'
import type { Meta, StoryObj } from '@storybook/react'
import { z } from 'zod'
import { renderDecorator } from '../utils/storybook'
import { DateInput } from './DateInput'

const label = 'Date'
const name = 'date'

const schema = z.object({
  date: z.date(),
})

const resolver = zodResolver(schema)

const defaultValues = {
  date: undefined,
}

const meta: Meta<typeof DateInput> = {
  title: 'Example/Form/DateInput',
  component: DateInput,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof DateInput>

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
