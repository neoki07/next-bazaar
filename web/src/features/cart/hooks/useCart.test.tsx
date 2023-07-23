import { setupMockServer } from '@/test-utils/mock-server'
import { createQueryWrapper } from '@/test-utils/wrappers'
import { QueryClient } from '@tanstack/react-query'
import { renderHook, waitFor } from '@testing-library/react'
import Decimal from 'decimal.js'
import { rest } from 'msw'
import { Cart } from '../types'
import { useCart } from './useCart'

const queryClient = new QueryClient()
const queryWrapper = createQueryWrapper(queryClient)

const handlers = [
  rest.get('*/cart', (_req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        products: [
          {
            id: '1',
            name: 'Product 1',
            description: 'Description 1',
            price: '10.0',
            quantity: 1,
            subtotal: '10.0',
            image_url: 'https://example.com/image.png',
          },
          {
            id: '2',
            name: 'Product 2',
            description: 'Description 2',
            price: '20.0',
            quantity: 2,
            subtotal: '40.0',
            image_url: 'https://example.com/image.png',
          },
        ],
        subtotal: '30.00',
        shipping: '5.00',
        tax: '3.00',
        total: '38.00',
      })
    )
  }),
]

describe('useCart', () => {
  setupMockServer(...handlers)

  it('returns cart correctly', async () => {
    const { result } = renderHook(() => useCart(), {
      wrapper: queryWrapper,
    })
    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    const expected: Cart = {
      products: [
        {
          id: '1',
          name: 'Product 1',
          description: 'Description 1',
          price: new Decimal(10.0),
          quantity: 1,
          subtotal: new Decimal(10.0),
          imageUrl: 'https://example.com/image.png',
        },
        {
          id: '2',
          name: 'Product 2',
          description: 'Description 2',
          price: new Decimal(20.0),
          quantity: 2,
          subtotal: new Decimal(40.0),
          imageUrl: 'https://example.com/image.png',
        },
      ],
      subtotal: new Decimal(30.0),
      shipping: new Decimal(5.0),
      tax: new Decimal(3.0),
      total: new Decimal(38.0),
    }
    expect(result.current.data).toEqual(expected)
  })
})
