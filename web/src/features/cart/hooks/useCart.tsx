import { useGetCart } from '@/api/endpoints/cart/cart'
import { CartDomainCartResponse } from '@/api/model'
import { AxiosResponse } from 'axios'
import { Cart } from '../types'
import { transformCart } from '../utils/transform'

function transform(response: AxiosResponse<CartDomainCartResponse>): Cart {
  const { data } = response
  return transformCart(data)
}

export function useCart() {
  return useGetCart({
    query: { select: transform },
    request: { withCredentials: true },
  })
}
