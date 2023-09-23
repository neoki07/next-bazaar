import { User } from '@/features/auth'

export const isTestUser = (user: User) => {
  return user.email === process.env.NEXT_PUBLIC_TEST_ACCOUNT_EMAIL
}
