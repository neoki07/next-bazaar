import { useGetProducts as useGetProductsQuery } from '@/api/endpoints/products/products'
import {
  ApiListProductsResponse,
  ApiListProductsResponseMeta,
} from '@/api/model'
import { transformProduct } from '@/features/products/utils/transform'
import { AxiosResponse } from 'axios'
import { Product } from '../types'

interface GetProductsResultData {
  meta: ApiListProductsResponseMeta
  data: Product[]
}

function transform(
  response: AxiosResponse<ApiListProductsResponse>
): GetProductsResultData {
  const { data } = response
  if (data.meta === undefined || data.data === undefined) {
    throw new Error('required fields are undefined:' + JSON.stringify(data))
  }

  return {
    meta: data.meta,
    data: data.data.map((item) => {
      return transformProduct(item)
    }),
  }
}

export function useGetProducts(page: number, pageSize: number) {
  return useGetProductsQuery<GetProductsResultData>(
    { page_id: page, page_size: pageSize },
    {
      query: { select: transform },
    }
  )
}
