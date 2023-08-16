import { TextInput, useForm } from '@/components/Form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button, Stack, Text, Title, rem } from '@mantine/core'
import { z } from 'zod'

export function PasswordSection() {
  const schema = z.object({
    oldPassword: z.string().min(8),
    newPassword: z.string().min(8),
    confirmNewPassword: z.string().min(8),
  })

  const [Form, methods] = useForm<{
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
      alert(JSON.stringify(data, null, 2))
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
          maw={rem(320)}
        />
        <TextInput
          label="New Password"
          name="newPassword"
          type="password"
          maw={rem(320)}
        />
        <TextInput
          label="Confirm New Password"
          name="confirmNewPassword"
          type="password"
          maw={rem(320)}
        />
      </Stack>
      <Text size="xs" mb={rem(4)}>
        {
          "Make sure it's at least 8 characters including a number and a letter."
        }
      </Text>
      <Button color="dark" type="submit">
        Update Password
      </Button>
    </Form>
  )
}
