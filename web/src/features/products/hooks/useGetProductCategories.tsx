import {
  getGetProductsCategoriesQueryKey,
  useGetProductsCategories,
} from '@/api/endpoints/products/products'
import {
  ProductDomainListProductCategoriesResponse,
  ProductDomainListProductCategoriesResponseMeta,
} from '@/api/model'
import { addNonCredentialsToQueryKey } from '@/utils/query-key'
import { AxiosResponse } from 'axios'
import { Category } from '../types'
import { transformProductCategory } from '../utils/transform'

interface GetProductCategoriesResultData {
  meta: ProductDomainListProductCategoriesResponseMeta
  data: Category[]
}

function transform(
  response: AxiosResponse<ProductDomainListProductCategoriesResponse>
): GetProductCategoriesResultData {
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
    data: data.data.map(transformProductCategory),
  }
}

export function useGetProductCategories(page: number, pageSize: number) {
  const params = { page_id: page, page_size: pageSize }
  const originalQueryKey = getGetProductsCategoriesQueryKey(params)

  return useGetProductsCategories<GetProductCategoriesResultData>(params, {
    query: {
      queryKey: addNonCredentialsToQueryKey(originalQueryKey),
      select: transform,
    },
  })
}
