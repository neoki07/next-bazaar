import { setupMockServer } from '@/test-utils/mock-server'
import { createQueryWrapper } from '@/test-utils/wrappers'
import { QueryClient } from '@tanstack/react-query'
import { renderHook, waitFor } from '@testing-library/react'
import { rest } from 'msw'
import { Category as ProductCategory } from '../types'
import { useGetProductCategories } from './useGetProductCategories'

const queryClient = new QueryClient()
const queryWrapper = createQueryWrapper(queryClient)

const handlers = [
  rest.get('*/products/categories', (_req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json({
        meta: {
          page_id: 1,
          page_size: 5,
        },
        data: [
          {
            id: '1',
            name: 'Category 1',
          },
          {
            id: '2',
            name: 'Category 2',
          },
          {
            id: '3',
            name: 'Category 3',
          },
          {
            id: '4',
            name: 'Category 4',
          },
          {
            id: '5',
            name: 'Category 5',
          },
        ],
      })
    )
  }),
]

describe('useGetProductCategories', () => {
  setupMockServer(...handlers)

  it('returns meta data and categories correctly', async () => {
    const { result } = renderHook(() => useGetProductCategories(1, 5), {
      wrapper: queryWrapper,
    })
    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    const expectedMeta = {
      page_id: 1,
      page_size: 5,
    }
    const expectedData: ProductCategory[] = [
      {
        id: '1',
        name: 'Category 1',
      },
      {
        id: '2',
        name: 'Category 2',
      },
      {
        id: '3',
        name: 'Category 3',
      },
      {
        id: '4',
        name: 'Category 4',
      },
      {
        id: '5',
        name: 'Category 5',
      },
    ]
    expect(result.current.data?.meta).toEqual(expectedMeta)
    expect(result.current.data?.data).toEqual(expectedData)
  })
})
