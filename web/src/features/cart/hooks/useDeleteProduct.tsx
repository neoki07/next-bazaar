import { ErrorType } from '@/api/custom-axios-instance'
import { useDeleteCartProducts } from '@/api/endpoints/cart/cart'
import { ApiErrorResponse, CartDomainDeleteProductRequest } from '@/api/model'
import { AxiosResponse } from 'axios'

interface UseDeleteProductParams {
  onSuccess?: (
    data: AxiosResponse<void, any>,
    variables: {
      data: CartDomainDeleteProductRequest
    },
    context: unknown
  ) => unknown
  onError?: (error: ErrorType<ApiErrorResponse>) => void
}

export function useDeleteProduct(params?: UseDeleteProductParams) {
  return useDeleteCartProducts({
    mutation: {
      onSuccess: params?.onSuccess,
      onError: params?.onError,
    },
    request: {
      withCredentials: true,
    },
  })
}
