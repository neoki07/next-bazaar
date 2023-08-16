import { TextInput, useForm } from '@/components/Form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button, Flex, rem } from '@mantine/core'
import { z } from 'zod'

interface NameSectionProps {
  initialName: string
}

export function NameSection({ initialName }: NameSectionProps) {
  const schema = z.object({
    name: z.string().min(1, { message: 'Required' }),
  })

  const [Form, methods] = useForm<{
    name: string
  }>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: initialName,
    },
    onSubmit: (data) => {
      alert(JSON.stringify(data, null, 2))
    },
  })

  return (
    <Form>
      <Flex gap="md" align="end" maw={rem(360)}>
        <div style={{ flex: 1 }}>
          <TextInput label="Username" name="name" />
        </div>
        <Button type="submit" color="dark">
          Save
        </Button>
      </Flex>
    </Form>
  )
}
