import { useGetProductsId as useGetProductQuery } from '@/api/endpoints/products/products'
import { ApiProductResponse } from '@/api/model'
import { transformProduct } from '@/features/products/utils/transform'
import { AxiosResponse } from 'axios'
import { Product } from '../types'

function transform(response: AxiosResponse<ApiProductResponse>): Product {
  return transformProduct(response.data)
}

export function useGetProduct(id: string) {
  return useGetProductQuery<Product>(id, {
    query: { select: transform },
  })
}
