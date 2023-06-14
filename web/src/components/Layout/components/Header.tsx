import { useAuth } from '@/features/auth'
import { useSession } from '@/providers/session'
import {
  Button,
  Group,
  Header as MantineHeader,
  Menu,
  Text,
  UnstyledButton,
  createStyles,
  rem,
} from '@mantine/core'
import {
  IconChevronDown,
  IconLogout,
  IconSettings,
  IconShoppingCart,
} from '@tabler/icons-react'
import Image from 'next/image'
import Link from 'next/link'
import { useCallback, useState } from 'react'

const useStyles = createStyles((theme) => ({
  user: {
    borderRadius: theme.radius.sm,
    transition: 'background-color 100ms ease',
  },
}))

export function Header() {
  const { session, status } = useSession()
  const { classes, theme } = useStyles()

  const [isLogoutButtonClicked, setIsLogoutButtonClicked] = useState(false)

  const handleLogoutError = useCallback(() => {
    setIsLogoutButtonClicked(false)
  }, [])

  const { logout } = useAuth({
    onLogoutError: handleLogoutError,
  })

  const handleLogout = useCallback(() => {
    setIsLogoutButtonClicked(true)
    logout()
  }, [logout])

  return (
    <MantineHeader height={60} px="md">
      <Group position="apart" sx={{ height: '100%' }}>
        <Link href="/">
          <Image src="/logo.svg" alt="Logo" width={136} height={36} priority />
        </Link>

        {status === 'authenticated' && (
          <Group spacing="xl" align="center">
            <Link
              href="/cart"
              style={{ display: 'flex', alignItems: 'center' }}
            >
              <IconShoppingCart size={24} stroke={1.5} />
            </Link>
            <Menu
              width={200}
              position="bottom-end"
              transitionProps={{ transition: 'pop-top-right' }}
              withinPortal
            >
              <Menu.Target>
                <UnstyledButton className={classes.user}>
                  <Group spacing={7}>
                    <Text weight={500} size="sm" sx={{ lineHeight: 1 }} mr={3}>
                      {session?.user.name}
                    </Text>
                    <IconChevronDown size={rem(12)} stroke={1.5} />
                  </Group>
                </UnstyledButton>
              </Menu.Target>
              <Menu.Dropdown>
                <Menu.Item icon={<IconSettings size="0.9rem" stroke={1.5} />}>
                  Account settings
                </Menu.Item>

                <Menu.Item
                  icon={<IconLogout size="0.9rem" stroke={1.5} />}
                  onClick={handleLogout}
                  disabled={isLogoutButtonClicked}
                >
                  Logout
                </Menu.Item>
              </Menu.Dropdown>
            </Menu>
          </Group>
        )}

        {status === 'unauthenticated' && (
          <Group>
            <Link href="/login">
              <Button variant="default">Log in</Button>
            </Link>
            <Link href="/register">
              <Button>Sign up</Button>
            </Link>
          </Group>
        )}
      </Group>
    </MantineHeader>
  )
}
