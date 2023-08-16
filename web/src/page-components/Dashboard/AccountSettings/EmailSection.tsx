import { TextInput, useForm } from '@/components/Form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button, Flex, rem } from '@mantine/core'
import { z } from 'zod'

interface EmailSectionProps {
  initialEmail: string
}

export function EmailSection({ initialEmail }: EmailSectionProps) {
  const schema = z.object({
    email: z.string().email({ message: 'Invalid email address' }),
  })

  const [Form, methods] = useForm<{
    email: string
  }>({
    resolver: zodResolver(schema),
    defaultValues: {
      email: initialEmail,
    },
    onSubmit: (data) => {
      alert(JSON.stringify(data, null, 2))
    },
  })

  return (
    <Form>
      <Flex gap="md" align="end" maw={rem(360)}>
        <div style={{ flex: 1 }}>
          <TextInput label="Email" name="email" />
        </div>
        <Button type="submit" color="dark">
          Save
        </Button>
      </Flex>
    </Form>
  )
}
