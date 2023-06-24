import { ErrorType } from '@/api/custom-axios-instance'
import { useDeleteCartProductId } from '@/api/endpoints/cart/cart'
import { ApiErrorResponse } from '@/api/model'
import { AxiosResponse } from 'axios'

interface UseDeleteProductParams {
  onSuccess?: (
    data: AxiosResponse<void, any>,
    variables: {
      productId: string
    },
    context: unknown
  ) => unknown
  onError?: (error: ErrorType<ApiErrorResponse>) => void
}

export function useDeleteProduct(params?: UseDeleteProductParams) {
  return useDeleteCartProductId({
    mutation: {
      onSuccess: params?.onSuccess,
      onError: params?.onError,
    },
    request: {
      withCredentials: true,
    },
  })
}
