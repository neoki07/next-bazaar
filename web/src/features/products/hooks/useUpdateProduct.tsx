import { ErrorType } from '@/api/custom-axios-instance'
import { usePutUsersProductsId } from '@/api/endpoints/users/users'
import {
  ApiErrorResponse,
  ApiMessageResponse,
  ProductDomainUpdateProductRequestBody,
} from '@/api/model'
import { AxiosResponse } from 'axios'

interface UseUpdateProductParams {
  onSuccess?: (
    data: AxiosResponse<ApiMessageResponse, any>,
    variables: {
      data: ProductDomainUpdateProductRequestBody
    },
    context: unknown
  ) => unknown
  onError?: (error: ErrorType<ApiErrorResponse>) => void
}

export function useUpdateProduct(params?: UseUpdateProductParams) {
  return usePutUsersProductsId({
    mutation: {
      onSuccess: params?.onSuccess,
      onError: params?.onError,
    },
    request: {
      withCredentials: true,
    },
  })
}
