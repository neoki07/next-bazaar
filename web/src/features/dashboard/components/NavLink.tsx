import { NavLink as MantineNavLink, createStyles } from '@mantine/core'
import { useRouter } from 'next/router'
import { ReactNode, useCallback } from 'react'

const useStyles = createStyles((theme) => ({
  link: {
    borderRadius: theme.radius.sm,

    '&:hover': {
      backgroundColor: theme.colors.gray[1],
    },
  },
}))

interface NavLinkProps {
  label: string
  icon: ReactNode
  pathname: string
  active: boolean
}

export function NavLink({ label, icon, pathname, active }: NavLinkProps) {
  const { classes } = useStyles()
  const router = useRouter()

  const handleClick = useCallback(() => {
    router.push(pathname)
  }, [router, pathname])

  return (
    <MantineNavLink
      className={classes.link}
      label={label}
      icon={icon}
      active={active}
      onClick={handleClick}
    />
  )
}
