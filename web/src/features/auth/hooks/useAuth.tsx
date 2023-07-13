import {
  usePostUsersLogin,
  usePostUsersLogout,
  usePostUsersRegister,
} from '@/api/endpoints/users/users'
import { useRouter } from 'next/router'
import { useCallback } from 'react'

interface LoginParams {
  email: string
  password: string
}

interface RegisterParams {
  name: string
  email: string
  password: string
}

interface UseAuthProps {
  onLoginError?: () => void
  onLogoutError?: () => void
  onRegisterError?: () => void
}

interface UseAuthResult {
  login: (params: LoginParams) => void
  logout: () => void
  registerAndLogin: (params: RegisterParams) => void
}

export function useAuth(props?: UseAuthProps): UseAuthResult {
  const router = useRouter()

  const loginMutation = usePostUsersLogin({
    mutation: {
      onSuccess: () => {
        router.push('/')
      },
      onError: () => {
        props?.onLoginError?.()
      },
    },
    request: {
      withCredentials: true,
    },
  })

  const logoutMutation = usePostUsersLogout({
    mutation: {
      onSuccess: () => {
        router.push('/')
      },
      onError: () => {
        props?.onLogoutError?.()
      },
    },
    request: {
      withCredentials: true,
    },
  })

  const registerMutation = usePostUsersRegister({
    mutation: {
      onError: () => {
        props?.onRegisterError?.()
      },
    },
    request: {
      withCredentials: true,
    },
  })

  const login = useCallback(
    ({ email, password }: LoginParams) => {
      loginMutation.mutate({ data: { email, password } })
    },
    [loginMutation]
  )

  const logout = useCallback(() => {
    logoutMutation.mutate()
  }, [logoutMutation])

  const registerAndLogin = useCallback(
    async ({ email, password, name }: RegisterParams) => {
      await registerMutation
        .mutateAsync({
          data: {
            email,
            password,
            name,
          },
        })
        .then(() => login({ email, password }))
        .catch((error: Error) => {
          console.error(error.message)
        })
    },
    [registerMutation, login]
  )

  return { login, logout, registerAndLogin }
}
