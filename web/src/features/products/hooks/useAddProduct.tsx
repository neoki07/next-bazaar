import { ErrorType } from '@/api/custom-axios-instance'
import { usePostUsersProducts } from '@/api/endpoints/users/users'
import {
  ApiErrorResponse,
  ApiMessageResponse,
  ProductDomainAddProductRequest,
} from '@/api/model'
import { AxiosResponse } from 'axios'

interface UseAddProductParams {
  onSuccess?: (
    data: AxiosResponse<ApiMessageResponse, any>,
    variables: {
      data: ProductDomainAddProductRequest
    },
    context: unknown
  ) => unknown
  onError?: (error: ErrorType<ApiErrorResponse>) => void
}

export function useAddProduct(params?: UseAddProductParams) {
  return usePostUsersProducts({
    mutation: {
      onSuccess: params?.onSuccess,
      onError: params?.onError,
    },
    request: {
      withCredentials: true,
    },
  })
}
