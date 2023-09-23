import { TextInput, useForm } from '@/components/Form'
import { useUpdateUserPassword } from '@/features/user'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button, Stack, Text, Title, rem } from '@mantine/core'
import { z } from 'zod'

interface PasswordSectionProps {
  isCurrentUserTestUser: boolean
  saving?: boolean
  onSubmit?: () => void
  onSubmitSuccess?: () => void
}

export function PasswordSection({
  isCurrentUserTestUser,
  saving,
  onSubmit,
  onSubmitSuccess,
}: PasswordSectionProps) {
  const schema = z
    .object({
      oldPassword: z.string().min(8),
      newPassword: z.string().min(8),
      confirmNewPassword: z.string().min(8),
    })
    .refine(
      ({ newPassword, confirmNewPassword }) => {
        return newPassword === confirmNewPassword
      },
      {
        path: ['confirmNewPassword'],
        message: 'Passwords do not match',
      }
    )
    .refine(
      ({ oldPassword, newPassword }) => {
        return oldPassword !== newPassword
      },
      {
        path: ['newPassword'],
        message: 'New password must be different from old password',
      }
    )

  const updatePasswordMutation = useUpdateUserPassword({
    onSuccess: onSubmitSuccess,
    onError: (error) => {
      throw new Error(error.message)
    },
  })

  const [Form] = useForm<{
    oldPassword: string
    newPassword: string
    confirmNewPassword: string
  }>({
    resolver: zodResolver(schema),
    defaultValues: {
      oldPassword: '',
      newPassword: '',
      confirmNewPassword: '',
    },
    onSubmit: (data) => {
      updatePasswordMutation.mutate({
        data: {
          old_password: data.oldPassword,
          new_password: data.newPassword,
        },
      })
      onSubmit?.()
    },
  })

  return (
    <Form>
      <Stack spacing="xs" mb="sm">
        <Title order={3}>Change Password</Title>
        <TextInput
          label="Old Password"
          name="oldPassword"
          type="password"
          maw={rem(360)}
          disabled={isCurrentUserTestUser}
        />
        <TextInput
          label="New Password"
          name="newPassword"
          type="password"
          maw={rem(360)}
          disabled={isCurrentUserTestUser}
        />
        <TextInput
          label="Confirm New Password"
          name="confirmNewPassword"
          type="password"
          maw={rem(360)}
          disabled={isCurrentUserTestUser}
        />
      </Stack>
      <Text size="xs" mb={rem(4)}>
        {"Make sure it's at least 8 characters."}
      </Text>
      <Button
        color="dark"
        type="submit"
        disabled={isCurrentUserTestUser || saving}
      >
        Update Password
      </Button>
    </Form>
  )
}
