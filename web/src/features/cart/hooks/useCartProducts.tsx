import { useGetCartProducts } from '@/api/endpoints/cart/cart'
import { CartDomainCartProductResponse } from '@/api/model'
import { AxiosResponse } from 'axios'
import { CartProduct } from '../types'
import { transformCartProduct } from '../utils/transform'

function transform(
  response: AxiosResponse<CartDomainCartProductResponse[]>
): CartProduct[] {
  const { data } = response
  if (data === undefined) {
    throw new Error('required fields are undefined:' + JSON.stringify(data))
  }

  return data.map((item) => transformCartProduct(item))
}

export function useCartProducts() {
  return useGetCartProducts({
    query: { select: transform },
    request: { withCredentials: true },
  })
}
