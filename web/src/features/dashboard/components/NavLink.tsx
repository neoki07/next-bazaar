import { NavLink as MantineNavLink } from '@mantine/core'
import { useRouter } from 'next/router'
import { ReactNode, useCallback } from 'react'

interface NavLinkProps {
  label: string
  icon: ReactNode
  pathname: string
  active: boolean
}

export function NavLink({ label, icon, pathname, active }: NavLinkProps) {
  const router = useRouter()

  const handleClick = useCallback(() => {
    router.push(pathname)
  }, [router, pathname])

  return (
    <MantineNavLink
      label={label}
      icon={icon}
      active={active}
      onClick={handleClick}
    />
  )
}
