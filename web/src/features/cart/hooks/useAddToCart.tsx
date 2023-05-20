import { ErrorType } from '@/api/custom-axios-instance'
import { usePostCartProducts } from '@/api/endpoints/cart/cart'
import { ApiErrorResponse } from '@/api/model'

interface UseAddToCartParams {
  onError?: (error: ErrorType<ApiErrorResponse>) => void
}

export function useAddToCart(params?: UseAddToCartParams) {
  return usePostCartProducts({
    mutation: {
      onError: params?.onError,
    },
    request: {
      withCredentials: true,
    },
  })
}
