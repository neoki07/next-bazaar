import { PasswordInput, TextInput, useForm } from '@/components/Form'
import { useAuth } from '@/features/auth'
import { zodResolver } from '@hookform/resolvers/zod'
import {
  Anchor,
  Button,
  Center,
  Container,
  Divider,
  Paper,
  Stack,
  Text,
  Title,
} from '@mantine/core'
import Image from 'next/image'
import Link from 'next/link'
import { useCallback, useState } from 'react'
import { z } from 'zod'

const schema = z.object({
  email: z
    .string()
    .min(1, { message: 'Required' })
    .email({ message: 'Wrong Format' }),
  password: z.string().min(8, { message: 'Minimum 8 characters' }),
})

export function Login() {
  const [isLoginButtonClicked, setIsLoginButtonClicked] = useState(false)

  const handleLoginError = useCallback(() => {
    setIsLoginButtonClicked(false)
  }, [])

  const { login } = useAuth({
    onLoginError: handleLoginError,
  })

  const loginWithTestAccount = useCallback(() => {
    setIsLoginButtonClicked(true)
    login({
      email: process.env.NEXT_PUBLIC_TEST_ACCOUNT_EMAIL || '',
      password: process.env.NEXT_PUBLIC_TEST_ACCOUNT_PASSWORD || '',
    })
  }, [login])

  const [Form, methods] = useForm<{
    email: string
    password: string
  }>({
    resolver: zodResolver(schema),
    defaultValues: {
      email: '',
      password: '',
    },
    onSubmit: (data) => {
      setIsLoginButtonClicked(true)
      login(data)
    },
  })

  return (
    <Container size={420} my={40}>
      <Stack>
        <Center>
          <Link href="/">
            <Image
              src="/logo.svg"
              alt="Logo"
              width={180}
              height={48}
              priority
            />
          </Link>
        </Center>

        <Title
          align="center"
          sx={{
            fontWeight: 900,
          }}
        >
          Welcome back!
        </Title>

        <Paper withBorder shadow="md" p={30} radius="md">
          <Form>
            <TextInput
              type="email"
              name="email"
              label="Email"
              placeholder="you@email.com"
            />
            <PasswordInput
              name="password"
              label="Password"
              placeholder="Your password"
              mt="md"
            />
            <Button
              type="submit"
              fullWidth
              mt="xl"
              disabled={isLoginButtonClicked}
            >
              Log in
            </Button>
          </Form>
        </Paper>

        <Text color="dimmed" size="sm" align="center" mt={5}>
          Do not have an account yet?{' '}
          <Link href="/register">
            <Anchor size="sm" component="button">
              Create account
            </Anchor>
          </Link>
        </Text>

        <Divider my="xs" label="OR" labelPosition="center" />

        <Button
          variant="outline"
          fullWidth
          disabled={isLoginButtonClicked}
          onClick={loginWithTestAccount}
        >
          Log in with Test Account
        </Button>
      </Stack>
    </Container>
  )
}
