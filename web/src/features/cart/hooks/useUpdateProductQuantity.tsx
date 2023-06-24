import { ErrorType } from '@/api/custom-axios-instance'
import { usePutCartProductId } from '@/api/endpoints/cart/cart'
import {
  ApiErrorResponse,
  ApiMessageResponse,
  CartDomainUpdateProductQuantityRequestBody,
} from '@/api/model'
import { AxiosResponse } from 'axios'

interface UseUpdateProductQuantityParams {
  onSuccess?: (
    data: AxiosResponse<ApiMessageResponse, any>,
    variables: {
      productId: string
      data: CartDomainUpdateProductQuantityRequestBody
    },
    context: unknown
  ) => unknown
  onError?: (error: ErrorType<ApiErrorResponse>) => void
}

export function useUpdateProductQuantity(
  params?: UseUpdateProductQuantityParams
) {
  return usePutCartProductId({
    mutation: {
      onSuccess: params?.onSuccess,
      onError: params?.onError,
    },
    request: {
      withCredentials: true,
    },
  })
}
