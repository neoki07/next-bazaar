import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { renderHook, waitFor } from '@testing-library/react'
import Decimal from 'decimal.js'
import { rest } from 'msw'
import { setupServer } from 'msw/node'
import { ReactNode } from 'react'
import { Product } from '../types'
import { useGetProduct } from './useGetProduct'

function setupGetProductsMockServer() {
  const server = setupServer(
    rest.get('*/products/:id', (_req, res, ctx) => {
      return res(
        ctx.delay(0),
        ctx.status(200, 'Mocked status'),
        ctx.json({
          id: '1',
          name: 'Product 1',
          description: 'Description 1',
          price: 10.0,
          stock_quantity: 10,
          category: 'Category 1',
          seller: 'Seller 1',
          image_url: 'https://example.com/image.png',
        })
      )
    })
  )

  beforeAll(() => server.listen())
  afterEach(() => server.resetHandlers())
  afterAll(() => server.close())

  return server
}

const queryCLient = new QueryClient()

const wrapper = ({ children }: { children: ReactNode }) => (
  <QueryClientProvider client={queryCLient}>{children}</QueryClientProvider>
)

describe('useGetProduct', () => {
  const server = setupGetProductsMockServer()

  it('returns product correctly', async () => {
    const { result } = renderHook(() => useGetProduct('1'), {
      wrapper,
    })
    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    const expected: Product = {
      id: '1',
      name: 'Product 1',
      description: 'Description 1',
      price: new Decimal(10.0),
      stockQuantity: 10,
      category: 'Category 1',
      seller: 'Seller 1',
      imageUrl: 'https://example.com/image.png',
    }
    expect(result.current.data).toEqual(expected)
  })
})
