import { TextInput, useForm } from '@/components/Form'
import { User } from '@/features/auth/types'
import { useUpdateUser } from '@/features/user'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button, Flex, rem } from '@mantine/core'
import { z } from 'zod'

interface EmailSectionProps {
  user: User
  disabledSaveButton?: boolean
  onSubmit?: () => void
  onSubmitSuccess?: () => void
}

export function EmailSection({
  user,
  disabledSaveButton,
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
          <TextInput label="Email" name="email" />
        </div>
        <Button type="submit" color="dark" disabled={disabledSaveButton}>
          Save
        </Button>
      </Flex>
    </Form>
  )
}
