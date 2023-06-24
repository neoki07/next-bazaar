import { useGetCartCount } from '@/api/endpoints/cart/cart'
import { CartDomainCartProductsCountResponse } from '@/api/model'
import { UseQueryOptions } from '@tanstack/react-query'
import { AxiosResponse } from 'axios'

function transform(
  response: AxiosResponse<CartDomainCartProductsCountResponse>
): number {
  const {
    data: { count },
  } = response

  if (count === undefined) {
    throw new Error('cart products count is undefined')
  }

  return count
}

export function useCartProductsCount(
  options?: UseQueryOptions<
    AxiosResponse<CartDomainCartProductsCountResponse, any>,
    unknown,
    number
  >
) {
  return useGetCartCount({
    query: { select: transform, ...options },
    request: { withCredentials: true },
  })
}
