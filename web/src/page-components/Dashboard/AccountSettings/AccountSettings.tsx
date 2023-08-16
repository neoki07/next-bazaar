import { useCurrentUser } from '@/features/auth'
import { DashboardLayout } from '@/features/dashboard'
import { Stack, Title, rem } from '@mantine/core'
import { useRouter } from 'next/router'
import { useCallback, useState } from 'react'
import { EmailSection } from './EmailSection'
import { NameSection } from './NameSection'
import { PasswordSection } from './PasswordSection'

export function AccountSettings() {
  const router = useRouter()
  const [saving, setSaving] = useState(false)
  const { data: user } = useCurrentUser()

  const handleSubmit = useCallback(() => {
    setSaving(true)
  }, [])

  const handleSubmitSuccess = useCallback(() => {
    router.reload()
  }, [router])

  return (
    <DashboardLayout>
      {user !== undefined && (
        <Stack spacing={rem(40)}>
          <Stack>
            <Title order={1}>Account Settings</Title>
            <NameSection
              user={user}
              disabledSaveButton={saving}
              onSubmit={handleSubmit}
              onSubmitSuccess={handleSubmitSuccess}
            />
            <EmailSection
              user={user}
              disabledSaveButton={saving}
              onSubmit={handleSubmit}
              onSubmitSuccess={handleSubmitSuccess}
            />
          </Stack>
          <PasswordSection disabledSaveButton={saving} />
        </Stack>
      )}
    </DashboardLayout>
  )
}
