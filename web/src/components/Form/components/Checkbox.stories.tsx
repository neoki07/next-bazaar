import { zodResolver } from '@hookform/resolvers/zod'
import type { Meta, StoryObj } from '@storybook/react'
import { z } from 'zod'
import { renderDecorator } from '../utils/storybook'
import { Checkbox } from './Checkbox'

const label = 'Checked'
const name = 'checked'

const schema = z.object({
  checked: z.boolean().refine((value) => value, { message: 'Required' }),
})

const resolver = zodResolver(schema)

const defaultValues = {
  checked: false,
}

const meta: Meta<typeof Checkbox> = {
  title: 'Example/Form/Checkbox',
  component: Checkbox,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof Checkbox>

export const Default: Story = {
  args: {
    label,
    name,
  },
  decorators: [(Story) => renderDecorator(Story, resolver, defaultValues)],
}

export const ErrorMessage: Story = {
  args: {
    label,
    name,
  },
  decorators: [
    (Story) => renderDecorator(Story, resolver, defaultValues, true),
  ],
}
