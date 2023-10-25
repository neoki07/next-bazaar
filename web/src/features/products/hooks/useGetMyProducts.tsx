import { useGetUsersProducts } from '@/api/endpoints/users/users'
import {
  ProductDomainListProductsResponse,
  ProductDomainListProductsResponseMeta,
} from '@/api/model'
import { transformProduct } from '@/features/products/utils/transform'
import { AxiosResponse } from 'axios'
import { Product } from '../types'

interface GetMyProductsResultData {
  meta: ProductDomainListProductsResponseMeta
  data: Product[]
}

function transform(
  response: AxiosResponse<ProductDomainListProductsResponse>
): GetMyProductsResultData {
  const { data } = response
  if (data.meta === undefined) {
    throw new Error(
      'required field `meta` is undefined:' + JSON.stringify(data)
    )
  } else if (data.data === undefined) {
    throw new Error(
      'required field `data` is undefined:' + JSON.stringify(data)
    )
  }

  return {
    meta: data.meta,
    data: data.data.map(transformProduct),
  }
}

interface UseGetMyProductsParams {
  page: number
  pageSize: number
}

export function useGetMyProducts({ page, pageSize }: UseGetMyProductsParams) {
  return useGetUsersProducts<GetMyProductsResultData>(
    {
      page_id: page,
      page_size: pageSize,
    },
    {
      query: { select: transform },
      request: { withCredentials: true },
    }
  )
}
