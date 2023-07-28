import { useGetProductsCategories } from '@/api/endpoints/products/products'
import {
  ProductDomainListProductCategoriesResponse,
  ProductDomainListProductCategoriesResponseMeta,
} from '@/api/model'
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

  console.log('data:', data)

  return {
    meta: data.meta,
    data: data.data.map(transformProductCategory),
  }
}

export function useGetProductCategories(page: number, pageSize: number) {
  return useGetProductsCategories<GetProductCategoriesResultData>(
    { page_id: page, page_size: pageSize },
    {
      query: { select: transform },
    }
  )
}
