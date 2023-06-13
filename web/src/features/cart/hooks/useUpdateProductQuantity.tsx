import { ErrorType } from '@/api/custom-axios-instance'
import { usePutCartProducts } from '@/api/endpoints/cart/cart'
import {
  ApiErrorResponse,
  ApiMessageResponse,
  CartDomainUpdateProductQuantityRequest,
} from '@/api/model'
import { AxiosResponse } from 'axios'

interface UseUpdateProductQuantityParams {
  onSuccess?: (
    data: AxiosResponse<ApiMessageResponse, any>,
    variables: {
      data: CartDomainUpdateProductQuantityRequest
    },
    context: unknown
  ) => unknown
  onError?: (error: ErrorType<ApiErrorResponse>) => void
}

export function useUpdateProductQuantity(
  params?: UseUpdateProductQuantityParams
) {
  return usePutCartProducts({
    mutation: {
      onSuccess: params?.onSuccess,
      onError: params?.onError,
    },
    request: {
      withCredentials: true,
    },
  })
}
