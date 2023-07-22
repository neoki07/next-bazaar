import { setupMockServer } from '@/test-utils/mock-server'
import { createQueryWrapper } from '@/test-utils/wrappers'
import { QueryClient } from '@tanstack/react-query'
import { renderHook, waitFor } from '@testing-library/react'
import { rest } from 'msw'
import { useCartProductsCount } from './useCartProductsCount'

const queryClient = new QueryClient()
const queryWrapper = createQueryWrapper(queryClient)

const handlers = [
  rest.get('*/cart/count', (_req, res, ctx) => {
    return res(ctx.status(200), ctx.json({ count: 5 }))
  }),
]

describe('useCartProductsCount', () => {
  setupMockServer(...handlers)

  it('returns cart products count correctly', async () => {
    const { result } = renderHook(() => useCartProductsCount(), {
      wrapper: queryWrapper,
    })
    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(result.current.data).toEqual(5)
  })
})
