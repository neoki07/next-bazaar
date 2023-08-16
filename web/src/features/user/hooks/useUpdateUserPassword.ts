import { ErrorType } from '@/api/custom-axios-instance'
import { usePatchUsersMePassword } from '@/api/endpoints/users/users'
import {
  ApiErrorResponse,
  ApiMessageResponse,
  UserDomainUpdatePasswordRequest,
} from '@/api/model'
import { AxiosResponse } from 'axios'

interface UseUpdateUserPasswordParams {
  onSuccess?: (
    data: AxiosResponse<ApiMessageResponse, any>,
    variables: {
      data: UserDomainUpdatePasswordRequest
    },
    context: unknown
  ) => unknown
  onError?: (error: ErrorType<ApiErrorResponse>) => void
}

export function useUpdateUserPassword(params?: UseUpdateUserPasswordParams) {
  return usePatchUsersMePassword({
    mutation: {
      onSuccess: params?.onSuccess,
      onError: params?.onError,
    },
    request: {
      withCredentials: true,
    },
  })
}
