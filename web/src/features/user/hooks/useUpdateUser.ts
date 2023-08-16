import { ErrorType } from '@/api/custom-axios-instance'
import { usePatchUsersMe } from '@/api/endpoints/users/users'
import {
  ApiErrorResponse,
  ApiMessageResponse,
  UserDomainUpdateRequest,
} from '@/api/model'
import { AxiosResponse } from 'axios'

interface UseUpdateUserParams {
  onSuccess?: (
    data: AxiosResponse<ApiMessageResponse, any>,
    variables: {
      data: UserDomainUpdateRequest
    },
    context: unknown
  ) => unknown
  onError?: (error: ErrorType<ApiErrorResponse>) => void
}

export function useUpdateUser(params?: UseUpdateUserParams) {
  return usePatchUsersMe({
    mutation: {
      onSuccess: params?.onSuccess,
      onError: params?.onError,
    },
    request: {
      withCredentials: true,
    },
  })
}
