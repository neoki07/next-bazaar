import { Meta, StoryObj } from '@storybook/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { CartProductInfoSkeleton } from './CartProductInfoSkeleton'

const queryClient = new QueryClient()

const meta: Meta<typeof CartProductInfoSkeleton> = {
  title: 'Features/Cart/CartProductInfoSkeleton',
  component: CartProductInfoSkeleton,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof CartProductInfoSkeleton>

export const Default: Story = {
  args: {
    imageSize: 300,
  },
  decorators: [
    (Story) => (
      <QueryClientProvider client={queryClient}>
        <Story />
      </QueryClientProvider>
    ),
  ],
}
