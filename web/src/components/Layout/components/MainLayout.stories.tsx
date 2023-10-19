import { Meta, StoryObj } from '@storybook/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MainLayout } from './MainLayout'

const queryClient = new QueryClient({})

const meta: Meta<typeof MainLayout> = {
  title: 'Components/Layout/MainLayout',
  component: MainLayout,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof MainLayout>

export const Default: Story = {
  args: {},
  decorators: [
    (Story) => (
      <QueryClientProvider client={queryClient}>
        <Story />
      </QueryClientProvider>
    ),
  ],
}
