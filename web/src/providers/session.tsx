import { AXIOS_INSTANCE } from '@/api/custom-axios-instance'
import { getCurrentUser } from '@/features/auth'
import { AxiosError } from 'axios'
import { useRouter } from 'next/router'
import {
  ReactNode,
  createContext,
  useContext,
  useEffect,
  useState,
} from 'react'

interface User {
  name: string
  email: string
}

export interface Session {
  user: User
}

export type SessionStatus = 'loading' | 'authenticated' | 'unauthenticated'

const SessionContext = createContext<{
  session: Session | undefined
  status: SessionStatus
}>({ session: undefined, status: 'loading' })

interface SessionProviderProps {
  children: ReactNode
}

export function SessionProvider({ children }: SessionProviderProps) {
  const router = useRouter()
  const [session, setSession] = useState<Session>()
  const [status, setStatus] = useState<SessionStatus>('loading')

  useEffect(() => {
    getCurrentUser()
      .then(({ data: { name, email } }) => {
        if (name !== undefined && email !== undefined) {
          setSession({ user: { name, email } })
          setStatus('authenticated')
        } else {
          throw new Error('username or email is undefined')
        }
      })
      .catch((error: AxiosError) => {
        if (error.response?.status === 401) {
          setSession(undefined)
          setStatus('unauthenticated')
        } else {
          throw new Error(error.message)
        }
      })
  }, [router])

  useEffect(() => {
    AXIOS_INSTANCE.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        if (error.response?.status === 401) {
          setSession(undefined)
          setStatus('unauthenticated')
        }

        return Promise.reject(error)
      }
    )
  }, [])

  return (
    <SessionContext.Provider value={{ session, status }}>
      {children}
    </SessionContext.Provider>
  )
}

export function useSession() {
  return useContext(SessionContext)
}
