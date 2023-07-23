import { setupMockServer } from '@/test-utils/mock-server'
import { createQueryWrapper } from '@/test-utils/wrappers'
import { QueryClient } from '@tanstack/react-query'
import { renderHook, waitFor } from '@testing-library/react'
import Decimal from 'decimal.js'
import { rest } from 'msw'
import { Product } from '../types'
import { useGetProduct } from './useGetProduct'

const queryClient = new QueryClient()
const queryWrapper = createQueryWrapper(queryClient)

const handlers = [
  rest.get('*/products/:id', (_req, res, ctx) => {
    return res(
      ctx.status(200),
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
  }),
]

describe('useGetProduct', () => {
  setupMockServer(...handlers)

  it('returns product correctly', async () => {
    const { result } = renderHook(() => useGetProduct('1'), {
      wrapper: queryWrapper,
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
