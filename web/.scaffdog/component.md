---
name: 'component'
root: 'src'
output: '**/*'
questions:
  name: 'Please enter component name.'
  component:
    confirm: 'Do you need component? (select no if you need only stories or tests for existing component)'
    initial: true
  stories:
    confirm: 'Do you need stories?'
    initial: true
  storiesTitle:
    if: inputs.stories
    message: '- Please enter stories title.'
  tests:
    confirm: 'Do you need tests?'
    initial: true
  testsUserEvent:
    if: inputs.tests
    confirm: '- Do you need user event in tests?'
    initial: true
  testsQuery:
    if: inputs.tests
    confirm: '- Do you need query client in tests?'
    initial: true
---

# `{{ !inputs.component && '!' }}{{ inputs.name | pascal }}.tsx`

```tsx
export interface {{ inputs.name | pascal }}Props {}

export function {{ inputs.name | pascal }}({}: {{ inputs.name | pascal }}Props) {
  return <div>This is {{ inputs.name | pascal }} component.</div>
}

```

# `{{ !inputs.stories && '!' }}{{ inputs.name | pascal }}.stories.tsx`

```tsx
import { Meta, StoryObj } from '@storybook/react'
import { {{ inputs.name | pascal }} } from './{{ inputs.name | pascal }}'

const meta: Meta<typeof {{ inputs.name | pascal }}> = {
  title: '{{ inputs.storiesTitle }}',
  component: {{ inputs.name | pascal }},
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof {{ inputs.name | pascal }}>

export const Default: Story = {
  args: {},
}

```

# `{{ !inputs.tests && '!' }}{{ inputs.name | pascal }}.test.tsx`

```tsx
{{ if inputs.testsQuery }}import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
{{ end }}import { render, screen } from '@testing-library/react'
{{ if inputs.testsUserEvent }}import userEvent from '@testing-library/user-event'
{{ end }}import { {{ inputs.name | pascal }} } from './{{ inputs.name | pascal }}'
{{ if inputs.testsUserEvent }}
const user = userEvent.setup()
{{ end }}{{ if inputs.testsQuery }}
const queryClient = new QueryClient({ logger: { ...console, error: () => {} } })
{{ end }}
describe('{{ inputs.name | pascal }}', () => {
  it('renders', () => {
    {{ if inputs.testsQuery }}render(
      <QueryClientProvider client={queryClient}>
        <{{ inputs.name | pascal }} />
      </QueryClientProvider>
    ){{ else }}render(<{{ inputs.name | pascal }} />){{ end }}

    expect(screen.getByText('This is {{ inputs.name | pascal }} component.')).toBeInTheDocument()
  })
})

```
