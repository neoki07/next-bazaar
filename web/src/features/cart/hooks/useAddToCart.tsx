import { ErrorType } from '@/api/custom-axios-instance'
import { usePostCartProducts } from '@/api/endpoints/cart/cart'
import {
  ApiErrorResponse,
  ApiMessageResponse,
  CartDomainAddProductRequest,
} from '@/api/model'
import { AxiosResponse } from 'axios'

interface UseAddToCartParams {
  onSuccess?: (
    data: AxiosResponse<ApiMessageResponse, any>,
    variables: {
      data: CartDomainAddProductRequest
    },
    context: unknown
  ) => unknown
  onError?: (error: ErrorType<ApiErrorResponse>) => void
}

export function useAddToCart(params?: UseAddToCartParams) {
  return usePostCartProducts({
    mutation: {
      onSuccess: params?.onSuccess,
      onError: params?.onError,
    },
    request: {
      withCredentials: true,
    },
  })
}
