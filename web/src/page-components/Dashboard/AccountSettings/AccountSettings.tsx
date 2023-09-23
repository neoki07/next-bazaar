import { useCurrentUser } from '@/features/auth'
import { DashboardLayout } from '@/features/dashboard'
import { isTestUser } from '@/features/user/utils/testuser'
import { Alert, Stack, Title, rem } from '@mantine/core'
import { IconInfoCircle } from '@tabler/icons-react'
import { useRouter } from 'next/router'
import { useCallback, useState } from 'react'
import { EmailSection } from './EmailSection'
import { NameSection } from './NameSection'
import { PasswordSection } from './PasswordSection'

export function AccountSettings() {
  const router = useRouter()
  const [saving, setSaving] = useState(false)
  const { data: user } = useCurrentUser()
  const isCurrentUserTestUser = user !== undefined && isTestUser(user)

  const handleSubmit = useCallback(() => {
    setSaving(true)
  }, [])

  const handleSubmitSuccess = useCallback(() => {
    router.reload()
  }, [router])

  return (
    <DashboardLayout>
      {user !== undefined && (
        <Stack>
          <Title order={1}>Account Settings</Title>
          {isCurrentUserTestUser && (
            <Alert
              variant="light"
              color="red"
              title="Account Information Cannot Be Changed"
              icon={<IconInfoCircle />}
              maw={rem(496)}
            >
              You are currently logged in as a test user. Please be aware that
              account information for test users cannot be changed.
            </Alert>
          )}
          <Stack mb="md">
            <NameSection
              user={user}
              isCurrentUserTestUser={isCurrentUserTestUser}
              saving={saving}
              onSubmit={handleSubmit}
              onSubmitSuccess={handleSubmitSuccess}
            />
            <EmailSection
              user={user}
              isCurrentUserTestUser={isCurrentUserTestUser}
              saving={saving}
              onSubmit={handleSubmit}
              onSubmitSuccess={handleSubmitSuccess}
            />
          </Stack>
          <PasswordSection
            isCurrentUserTestUser={isCurrentUserTestUser}
            saving={saving}
            onSubmit={handleSubmit}
            onSubmitSuccess={handleSubmitSuccess}
          />
        </Stack>
      )}
    </DashboardLayout>
  )
}
