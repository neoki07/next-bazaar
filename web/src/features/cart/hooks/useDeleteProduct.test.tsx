import { setupMockServer } from '@/test-utils/mock-server'
import { createQueryWrapper } from '@/test-utils/wrappers'
import { QueryClient } from '@tanstack/react-query'
import { renderHook, waitFor } from '@testing-library/react'
import { rest } from 'msw'
import { useDeleteProduct } from './useDeleteProduct'

const queryClient = new QueryClient({ logger: { ...console, error: () => {} } })
const queryWrapper = createQueryWrapper(queryClient)

const handlers = [
  rest.delete('*/cart/1', (_req, res, ctx) => {
    return res(ctx.status(204))
  }),
]

describe('useDeleteProduct', () => {
  const server = setupMockServer(...handlers)

  it('calls onSuccess callback when the request succeeds', async () => {
    const mockOnSuccess = jest.fn()
    const { result } = renderHook(
      () => useDeleteProduct({ onSuccess: mockOnSuccess }),
      { wrapper: queryWrapper }
    )

    result.current.mutate({ productId: '1' })

    await waitFor(() => {
      expect(mockOnSuccess).toHaveBeenCalledTimes(1)
    })
  })

  it('calls onError callback when the request fails', async () => {
    server.use(
      rest.delete('*/cart/1', (_req, res, ctx) => {
        return res(ctx.status(500))
      })
    )
    const mockOnError = jest.fn()
    const { result } = renderHook(
      () => useDeleteProduct({ onError: mockOnError }),
      { wrapper: queryWrapper }
    )

    result.current.mutate({ productId: '1' })

    await waitFor(() => {
      expect(mockOnError).toHaveBeenCalledTimes(1)
    })
  })
})
