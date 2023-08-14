import { Box, NavLink } from '@mantine/core'
import { IconPackageExport, IconSettings } from '@tabler/icons-react'

interface NavBarProps {
  width: number
}

export function NavBar({ width }: NavBarProps) {
  return (
    <Box w={width}>
      <NavLink
        label="Account Settings"
        icon={<IconSettings size="1rem" stroke={1.5} />}
      />
      <NavLink
        label="Your Products"
        icon={<IconPackageExport size="1rem" stroke={1.5} />}
        active
      />
    </Box>
  )
}
