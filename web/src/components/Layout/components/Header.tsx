import { useState } from "react";
import {
  createStyles,
  Avatar,
  UnstyledButton,
  Group,
  Text,
  Menu,
  rem,
} from "@mantine/core";
import { IconLogout, IconSettings, IconChevronDown } from "@tabler/icons-react";
import { MantineLogo } from "@mantine/ds";

const useStyles = createStyles((theme) => ({
  header: {
    padding: theme.spacing.md,
    backgroundColor: theme.fn.variant({
      variant: "filled",
      color: theme.primaryColor,
    }).background,
    borderBottom: `${rem(1)} solid ${
      theme.fn.variant({ variant: "filled", color: theme.primaryColor })
        .background
    }`,
  },

  user: {
    color: theme.white,
    // padding: `${theme.spacing.xs} ${theme.spacing.sm}`,
    borderRadius: theme.radius.sm,
    transition: "background-color 100ms ease",
  },
}));

interface HeaderProps {
  user: { name: string; image: string };
}

export function Header({ user }: HeaderProps) {
  const { classes, theme, cx } = useStyles();

  return (
    <div className={classes.header}>
      <Group position="apart">
        <MantineLogo size={28} inverted />

        <Menu
          width={200}
          position="bottom-end"
          transitionProps={{ transition: "pop-top-right" }}
          withinPortal
        >
          <Menu.Target>
            <UnstyledButton
              className={cx(classes.user, {
                // [classes.userActive]: userMenuOpened,
              })}
            >
              <Group spacing={7}>
                <Avatar
                  src={user.image}
                  alt={user.name}
                  radius="xl"
                  size={20}
                />
                <Text
                  weight={500}
                  size="sm"
                  sx={{ lineHeight: 1, color: theme.white }}
                  mr={3}
                >
                  {user.name}
                </Text>
                <IconChevronDown size={rem(12)} stroke={1.5} />
              </Group>
            </UnstyledButton>
          </Menu.Target>
          <Menu.Dropdown>
            <Menu.Item icon={<IconSettings size="0.9rem" stroke={1.5} />}>
              Account settings
            </Menu.Item>

            <Menu.Item icon={<IconLogout size="0.9rem" stroke={1.5} />}>
              Logout
            </Menu.Item>
          </Menu.Dropdown>
        </Menu>
      </Group>
    </div>
  );
}
