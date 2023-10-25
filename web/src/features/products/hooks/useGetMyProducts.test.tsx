import { setupMockServer } from '@/test-utils/mock-server'
import { createQueryWrapper } from '@/test-utils/wrappers'
import { QueryClient } from '@tanstack/react-query'
import { renderHook, waitFor } from '@testing-library/react'
import Decimal from 'decimal.js'
import { rest } from 'msw'
import { Product } from '../types'
import { useGetMyProducts } from './useGetMyProducts'

const queryClient = new QueryClient()
const queryWrapper = createQueryWrapper(queryClient)

const handlers = [
  rest.get('*/users/products', (_req, res, ctx) => {
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
            category_id: '1',
            category: 'Category 1',
            seller: 'Seller 1',
            image_url: 'https://example.com/image.png',
          },
          {
            id: '2',
            name: 'Product 2',
            description: 'Description 2',
            price: 20.0,
            stock_quantity: 20,
            category_id: '2',
            category: 'Category 2',
            seller: 'Seller 1',
            image_url: 'https://example.com/image.png',
          },
        ],
      })
    )
  }),
]

describe('useGetMyProducts', () => {
  setupMockServer(...handlers)

  it('returns meta data and products correctly', async () => {
    const { result } = renderHook(
      () =>
        useGetMyProducts({
          page: 1,
          pageSize: 5,
        }),
      {
        wrapper: queryWrapper,
      }
    )
    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    const expectedMeta = {
      page_id: 1,
      page_size: 5,
      total_count: 5,
      total_pages: 1,
    }
    const expectedData: Product[] = [
      {
        id: '1',
        name: 'Product 1',
        description: 'Description 1',
        price: new Decimal(10.0),
        stockQuantity: 10,
        categoryId: '1',
        category: 'Category 1',
        seller: 'Seller 1',
        imageUrl: 'https://example.com/image.png',
      },
      {
        id: '2',
        name: 'Product 2',
        description: 'Description 2',
        price: new Decimal(20.0),
        stockQuantity: 20,
        categoryId: '2',
        category: 'Category 2',
        seller: 'Seller 1',
        imageUrl: 'https://example.com/image.png',
      },
    ]
    expect(result.current.data?.meta).toEqual(expectedMeta)
    expect(result.current.data?.data).toEqual(expectedData)
  })
})
