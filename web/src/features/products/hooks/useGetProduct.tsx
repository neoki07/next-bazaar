import {
  getGetProductsIdQueryKey,
  useGetProductsId as useGetProductQuery,
} from '@/api/endpoints/products/products'
import { ProductDomainProductResponse } from '@/api/model'
import { transformProduct } from '@/features/products/utils/transform'
import { addNonCredentialsToQueryKey } from '@/utils/query-key'
import { AxiosResponse } from 'axios'
import { Product } from '../types'

function transform(
  response: AxiosResponse<ProductDomainProductResponse>
): Product {
  return transformProduct(response.data)
}

export function useGetProduct(id: string) {
  const originalQueryKey = getGetProductsIdQueryKey(id)

  return useGetProductQuery<Product>(id, {
    query: {
      queryKey: addNonCredentialsToQueryKey(originalQueryKey),
      select: transform,
    },
  })
}
