import { setupMockServer } from '@/test-utils/mock-server'
import { createQueryWrapper } from '@/test-utils/wrappers'
import { QueryClient } from '@tanstack/react-query'
import { renderHook, waitFor } from '@testing-library/react'
import { rest } from 'msw'
import { useAddProduct } from './useAddProduct'

const queryClient = new QueryClient({ logger: { ...console, error: () => {} } })
const queryWrapper = createQueryWrapper(queryClient)

const handlers = [
  rest.post('*/users/products', (_req, res, ctx) => {
    return res(ctx.status(200), ctx.json({ message: 'Mock' }))
  }),
]

describe('useAddProduct', () => {
  const server = setupMockServer(...handlers)

  it('calls onSuccess callback when the request succeeds', async () => {
    const mockOnSuccess = jest.fn()
    const { result } = renderHook(
      () => useAddProduct({ onSuccess: mockOnSuccess }),
      { wrapper: queryWrapper }
    )

    result.current.mutate({
      data: {
        name: 'Product 1',
        description: 'Description 1',
        price: '10.00',
        stock_quantity: 10,
        category_id: '1',
        image_url: 'https://example.com/image.png',
      },
    })

    await waitFor(() => {
      expect(mockOnSuccess).toHaveBeenCalledTimes(1)
    })
  })

  it('calls onError callback when the request fails', async () => {
    server.use(
      rest.post('*/users/products', (_req, res, ctx) => {
        return res(ctx.status(500))
      })
    )
    const mockOnError = jest.fn()
    const { result } = renderHook(
      () => useAddProduct({ onError: mockOnError }),
      { wrapper: queryWrapper }
    )

    result.current.mutate({
      data: {
        name: 'Product 1',
        description: 'Description 1',
        price: '10.00',
        stock_quantity: 10,
        category_id: '1',
        image_url: 'https://example.com/image.png',
      },
    })

    await waitFor(() => {
      expect(mockOnError).toHaveBeenCalledTimes(1)
    })
  })
})
