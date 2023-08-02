import { Meta, StoryObj } from '@storybook/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { CategorySectionSkeleton } from './CategorySectionSkeleton'

const IMAGE_SIZE = 260

const queryClient = new QueryClient({})

const meta: Meta<typeof CategorySectionSkeleton> = {
  title: 'PageComponents/Home/CategorySectionSkeleton',
  component: CategorySectionSkeleton,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof CategorySectionSkeleton>

export const Default: Story = {
  args: {
    imageSize: IMAGE_SIZE,
  },
  decorators: [
    (Story) => (
      <QueryClientProvider client={queryClient}>
        <Story />
      </QueryClientProvider>
    ),
  ],
}
