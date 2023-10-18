import { zodResolver } from '@hookform/resolvers/zod'
import type { Meta, StoryObj } from '@storybook/react'
import { z } from 'zod'
import { renderDecorator } from '../utils/storybook'
import { MultiSelect } from './MultiSelect'

const label = 'Programming Language'
const name = 'programmingLanguage'
const options = [
  {
    label: 'JavaScript',
    value: 'javascript',
  },
  {
    label: 'TypeScript',
    value: 'typescript',
  },
  {
    label: 'Go',
    value: 'go',
  },
  {
    label: 'Python',
    value: 'python',
  },
  {
    label: 'Rust',
    value: 'rust',
  },
]

const schema = z.object({
  programmingLanguage: z.string().min(1, { message: 'Required' }),
})

const resolver = zodResolver(schema)

const defaultValues = {
  programmingLanguage: [],
}

const meta: Meta<typeof MultiSelect> = {
  title: 'Components/Form/MultiSelect',
  component: MultiSelect,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof MultiSelect>

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
