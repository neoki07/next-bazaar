import { createStyles, UnstyledButton, Text, Menu, rem } from "@mantine/core";
import { IconLogout, IconSettings, IconChevronDown } from "@tabler/icons-react";
import { useSession } from "@/providers/session";
import { Header as MantineHeader, Group, Button, Box } from "@mantine/core";
import { MantineLogo } from "@mantine/ds";
import Link from "next/link";
import { useCallback, useState } from "react";
import { useAuth } from "@/features/auth";

const useStyles = createStyles((theme) => ({
  user: {
    borderRadius: theme.radius.sm,
    transition: "background-color 100ms ease",
  },
}));

export function Header() {
  const { session, status } = useSession();
  const { classes, theme } = useStyles();

  const [isLogoutButtonClicked, setIsLogoutButtonClicked] = useState(false);

  const handleLogoutError = useCallback(() => {
    setIsLogoutButtonClicked(false);
  }, []);

  const { logout } = useAuth({
    onLogoutError: handleLogoutError,
  });

  const handleLogout = useCallback(() => {
    setIsLogoutButtonClicked(true);
    logout();
  }, [logout]);

  return (
    <MantineHeader height={60} px="md">
      <Group position="apart" sx={{ height: "100%" }}>
        <MantineLogo size={30} />

        {status === "authenticated" && (
          <Menu
            width={200}
            position="bottom-end"
            transitionProps={{ transition: "pop-top-right" }}
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
        )}

        {status === "unauthenticated" && (
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
  );
}
