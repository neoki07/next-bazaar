import { useGetUsersMe } from '@/api/endpoints/users/users'
import { UserDomainUserResponse } from '@/api/model'
import { AxiosResponse } from 'axios'
import { User } from '../types'

function transform(response: AxiosResponse<UserDomainUserResponse>): User {
  const { data } = response

  if (data.name === undefined) {
    throw new Error(
      'required field `name` is undefined:' + JSON.stringify(data)
    )
  } else if (data.email === undefined) {
    throw new Error(
      'required field `email` is undefined:' + JSON.stringify(data)
    )
  }

  return {
    name: data.name,
    email: data.email,
  }
}

export function useCurrentUser() {
  return useGetUsersMe({
    query: {
      select: transform,
    },
    request: {
      withCredentials: true,
    },
  })
}
