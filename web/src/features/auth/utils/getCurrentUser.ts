import { getUsersMe } from '@/api/endpoints/users/users'

export function getCurrentUser() {
  return getUsersMe({ withCredentials: true })
}
