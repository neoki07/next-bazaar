import { setupMockServer } from '@/test-utils/mock-server'
import { createQueryWrapper } from '@/test-utils/wrappers'
import { QueryClient } from '@tanstack/react-query'
import { renderHook, waitFor } from '@testing-library/react'
import { rest } from 'msw'
import { useAddToCart } from './useAddToCart'

const queryClient = new QueryClient({ logger: { ...console, error: () => {} } })
const queryWrapper = createQueryWrapper(queryClient)

const handlers = [
  rest.post('*/cart/add-product', (_req, res, ctx) => {
    return res(ctx.status(200), ctx.json({ message: 'Mock' }))
  }),
]

describe('useAddToCart', () => {
  const server = setupMockServer(...handlers)

  it('calls onSuccess callback when the request succeeds', async () => {
    const mockOnSuccess = jest.fn()
    const { result } = renderHook(
      () => useAddToCart({ onSuccess: mockOnSuccess }),
      { wrapper: queryWrapper }
    )

    result.current.mutate({ data: { product_id: '1', quantity: 5 } })

    await waitFor(() => {
      expect(mockOnSuccess).toHaveBeenCalledTimes(1)
    })
  })

  it('calls onError callback when the request fails', async () => {
    server.use(
      rest.post('*/cart/add-product', (_req, res, ctx) => {
        return res(ctx.status(500))
      })
    )
    const mockOnError = jest.fn()
    const { result } = renderHook(
      () => useAddToCart({ onError: mockOnError }),
      { wrapper: queryWrapper }
    )

    result.current.mutate({ data: { product_id: '1', quantity: 5 } })

    await waitFor(() => {
      expect(mockOnError).toHaveBeenCalledTimes(1)
    })
  })
})
