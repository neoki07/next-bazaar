import { Meta, StoryObj } from '@storybook/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { rest } from 'msw'
import { CategorySection } from './CategorySection'

const IMAGE_SIZE = 260

const handlers = [
  rest.get('*/products', (_req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        meta: {
          page_id: 1,
          page_size: 5,
          total_count: 5,
          total_pages: 1,
        },
        data: [
          {
            id: '1',
            name: 'Product 1',
            description: 'Description 1',
            price: 10.0,
            stock_quantity: 10,
            category: 'Category 1',
            seller: 'Seller 1',
            image_url: 'https://via.placeholder.com/260',
          },
          {
            id: '2',
            name: 'Product 2',
            description: 'Description 2',
            price: 20.0,
            stock_quantity: 20,
            category: 'Category 2',
            seller: 'Seller 2',
            image_url: 'https://via.placeholder.com/260',
          },
          {
            id: '3',
            name: 'Product 3',
            description: 'Description 3',
            price: 30.0,
            stock_quantity: 30,
            category: 'Category 3',
            seller: 'Seller 3',
            image_url: 'https://via.placeholder.com/260',
          },
          {
            id: '4',
            name: 'Product 4',
            description: 'Description 4',
            price: 40.0,
            stock_quantity: 40,
            category: 'Category 4',
            seller: 'Seller 4',
            image_url: 'https://via.placeholder.com/260',
          },
          {
            id: '5',
            name: 'Product 5',
            description: 'Description 5',
            price: 50.0,
            stock_quantity: 50,
            category: 'Category 5',
            seller: 'Seller 5',
            image_url: 'https://via.placeholder.com/260',
          },
        ],
      })
    )
  }),
]

const queryClient = new QueryClient({})

const meta: Meta<typeof CategorySection> = {
  title: 'PageComponents/Home/CategorySection',
  component: CategorySection,
  tags: ['autodocs'],
  argTypes: {},
}

export default meta
type Story = StoryObj<typeof CategorySection>

export const Default: Story = {
  args: {
    category: {
      id: '1',
      name: 'Category 1',
    },
    getProductLink: () => '#',
  },
  decorators: [
    (Story) => (
      <QueryClientProvider client={queryClient}>
        <Story />
      </QueryClientProvider>
    ),
  ],
  parameters: {
    msw: {
      handlers,
    },
  },
}
