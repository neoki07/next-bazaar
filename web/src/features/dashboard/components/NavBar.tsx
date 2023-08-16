import { Box } from '@mantine/core'
import { links, useActiveLink } from '../links'
import { NavLink } from './NavLink'

interface NavBarProps {
  width: number
}

export function NavBar({ width }: NavBarProps) {
  const activeLink = useActiveLink()

  return (
    <Box w={width}>
      {links.map((link) => (
        <NavLink
          key={link.pathname}
          label={link.label}
          pathname={link.pathname}
          icon={link.icon}
          active={link === activeLink}
        />
      ))}
    </Box>
  )
}
