import { TextInput, useForm } from '@/components/Form'
import { User } from '@/features/auth'
import { useUpdateUser } from '@/features/user'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button, Flex, rem } from '@mantine/core'
import { z } from 'zod'

interface EmailSectionProps {
  user: User
  isCurrentUserTestUser: boolean
  saving?: boolean
  onSubmit?: () => void
  onSubmitSuccess?: () => void
}

export function EmailSection({
  user,
  isCurrentUserTestUser,
  saving,
  onSubmit,
  onSubmitSuccess,
}: EmailSectionProps) {
  const schema = z.object({
    email: z.string().email({ message: 'Invalid email address' }),
  })

  const updateUserMutation = useUpdateUser({
    onSuccess: onSubmitSuccess,
    onError: (error) => {
      throw new Error(error.message)
    },
  })

  const [Form] = useForm<{
    email: string
  }>({
    resolver: zodResolver(schema),
    defaultValues: {
      email: user.email,
    },
    onSubmit: (data) => {
      updateUserMutation.mutate({ data: { ...user, ...data } })
      onSubmit?.()
    },
  })

  return (
    <Form>
      <Flex gap="md" align="end" maw={rem(392)}>
        <div style={{ flex: 1 }}>
          <TextInput
            label="Email"
            name="email"
            disabled={isCurrentUserTestUser}
          />
        </div>
        <Button
          type="submit"
          color="dark"
          disabled={isCurrentUserTestUser || saving}
        >
          Save
        </Button>
      </Flex>
    </Form>
  )
}
