import { ReactNode, createContext, useContext } from 'react'

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
  session?: Session
  status: SessionStatus
  children: ReactNode
}

export function SessionProvider({
  session,
  status,
  children,
}: SessionProviderProps) {
  return (
    <SessionContext.Provider value={{ session, status }}>
      {children}
    </SessionContext.Provider>
  )
}

export function useSession() {
  return useContext(SessionContext)
}
