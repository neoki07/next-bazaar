import { useGetCart } from '@/api/endpoints/cart/cart'
import { CartDomainCartResponse } from '@/api/model'
import { AxiosResponse } from 'axios'
import { Cart } from '../types'
import { transformCart } from '../utils/transform'

function transform(response: AxiosResponse<CartDomainCartResponse>): Cart {
  const { data } = response
  if (data === undefined) {
    throw new Error('required fields are undefined:' + JSON.stringify(data))
  }

  return transformCart(data)
}

export function useCart() {
  return useGetCart({
    query: { select: transform },
    request: { withCredentials: true },
  })
}
