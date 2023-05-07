import {
  Anchor,
  Paper,
  Title,
  Text,
  Container,
  Button,
  Center,
  Stack,
} from "@mantine/core";
import { MantineLogo } from "@mantine/ds";
import Link from "next/link";
import { z } from "zod";
import { useCallback, useState } from "react";
import { useAuth } from "@/features/auth";
import { PasswordInput, TextInput, useForm } from "@/components/Form";
import { zodResolver } from "@hookform/resolvers/zod";

const schema = z
  .object({
    name: z.string().min(1, { message: "Required" }),
    email: z
      .string()
      .min(1, { message: "Required" })
      .email({ message: "Wrong Format" }),
    password: z.string().min(8),
    confirmPassword: z.string().min(8),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords do not match",
    path: ["confirmPassword"],
  });

export default function RegisterPage() {
  const [isRegisterButtonClicked, setIsRegisterButtonClicked] = useState(false);

  const handleRegisterError = useCallback(() => {
    setIsRegisterButtonClicked(false);
  }, []);

  const { registerAndLogin } = useAuth({
    onRegisterError: handleRegisterError,
  });

  const [Form, methods] = useForm<{
    name: string;
    email: string;
    password: string;
    confirmPassword: string;
  }>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: "",
      email: "",
      password: "",
      confirmPassword: "",
    },
    onSubmit: ({ confirmPassword: _, ...data }) => {
      setIsRegisterButtonClicked(true);
      registerAndLogin(data);
    },
  });

  return (
    <Container size={420} my={40}>
      <Stack>
        <Center>
          <Link href={"/"}>
            <MantineLogo size={30} />
          </Link>
        </Center>

        <Title
          align="center"
          sx={{
            fontWeight: 900,
          }}
        >
          Create new account
        </Title>

        <Paper withBorder shadow="md" p={30} radius="md">
          <Form>
            <TextInput name="name" label="Username" placeholder="Your name" />
            <TextInput
              type="email"
              name="email"
              label="Email"
              placeholder="you@email.com"
              mt="md"
            />
            <PasswordInput
              name="password"
              label="Password"
              placeholder="Your password"
              mt="md"
            />
            <PasswordInput
              name="confirmPassword"
              label="Confirm Password"
              placeholder="Your password"
              mt="md"
            />
            <Button
              type="submit"
              fullWidth
              mt="xl"
              disabled={isRegisterButtonClicked}
            >
              Sign up
            </Button>
          </Form>
        </Paper>

        <Text color="dimmed" size="sm" align="center" mt={5}>
          Already have an account?{" "}
          <Link href="/login">
            <Anchor size="sm" component="button">
              Login
            </Anchor>
          </Link>
        </Text>
      </Stack>
    </Container>
  );
}
