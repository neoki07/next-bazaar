import { useCurrentUser } from '@/features/auth'
import { DashboardLayout } from '@/features/dashboard'
import { Stack, Title, rem } from '@mantine/core'
import { EmailSection } from './EmailSection'
import { NameSection } from './NameSection'
import { PasswordSection } from './PasswordSection'

export function AccountSettings() {
  const { data: user } = useCurrentUser()

  return (
    <DashboardLayout>
      {user !== undefined && (
        <Stack spacing={rem(40)}>
          <Stack>
            <Title order={1}>Account Settings</Title>
            <NameSection initialName={user.name} />
            <EmailSection initialEmail={user.email} />
          </Stack>
          <PasswordSection />
        </Stack>
      )}
    </DashboardLayout>
  )
}
