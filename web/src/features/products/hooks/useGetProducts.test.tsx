import { setupMockServer } from '@/test-utils/mock-server'
import { createQueryWrapper } from '@/test-utils/wrappers'
import { QueryClient } from '@tanstack/react-query'
import { renderHook, waitFor } from '@testing-library/react'
import Decimal from 'decimal.js'
import { rest } from 'msw'
import { Product } from '../types'
import { useGetProducts } from './useGetProducts'

const queryClient = new QueryClient()
const queryWrapper = createQueryWrapper(queryClient)

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
            image_url: 'https://example.com/image.png',
          },
          {
            id: '2',
            name: 'Product 2',
            description: 'Description 2',
            price: 20.0,
            stock_quantity: 20,
            category: 'Category 2',
            seller: 'Seller 2',
            image_url: 'https://example.com/image.png',
          },
          {
            id: '3',
            name: 'Product 3',
            description: 'Description 3',
            price: 30.0,
            stock_quantity: 30,
            category: 'Category 3',
            seller: 'Seller 3',
            image_url: 'https://example.com/image.png',
          },
          {
            id: '4',
            name: 'Product 4',
            description: 'Description 4',
            price: 40.0,
            stock_quantity: 40,
            category: 'Category 4',
            seller: 'Seller 4',
            image_url: 'https://example.com/image.png',
          },
          {
            id: '5',
            name: 'Product 5',
            description: 'Description 5',
            price: 50.0,
            stock_quantity: 50,
            category: 'Category 5',
            seller: 'Seller 5',
            image_url: 'https://example.com/image.png',
          },
        ],
      })
    )
  }),
]

describe('useGetProducts', () => {
  setupMockServer(...handlers)

  it('returns meta data and products correctly', async () => {
    const { result } = renderHook(
      () =>
        useGetProducts({
          page: 1,
          pageSize: 5,
          categoryId: 'dummy',
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
        category: 'Category 2',
        seller: 'Seller 2',
        imageUrl: 'https://example.com/image.png',
      },
      {
        id: '3',
        name: 'Product 3',
        description: 'Description 3',
        price: new Decimal(30.0),
        stockQuantity: 30,
        category: 'Category 3',
        seller: 'Seller 3',
        imageUrl: 'https://example.com/image.png',
      },
      {
        id: '4',
        name: 'Product 4',
        description: 'Description 4',
        price: new Decimal(40.0),
        stockQuantity: 40,
        category: 'Category 4',
        seller: 'Seller 4',
        imageUrl: 'https://example.com/image.png',
      },
      {
        id: '5',
        name: 'Product 5',
        description: 'Description 5',
        price: new Decimal(50.0),
        stockQuantity: 50,
        category: 'Category 5',
        seller: 'Seller 5',
        imageUrl: 'https://example.com/image.png',
      },
    ]
    expect(result.current.data?.meta).toEqual(expectedMeta)
    expect(result.current.data?.data).toEqual(expectedData)
  })
})
