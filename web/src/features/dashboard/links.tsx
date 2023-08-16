import { IconPackageExport, IconSettings } from '@tabler/icons-react'
import { useRouter } from 'next/router'
import { ReactNode } from 'react'

interface Link {
  label: string
  icon: ReactNode
  pathname: string
}

export const links: Link[] = [
  {
    label: 'Account Settings',
    icon: <IconSettings size="1rem" stroke={1.5} />,
    pathname: '/dashboard/settings/account',
  },
  {
    label: 'Your Products',
    icon: <IconPackageExport size="1rem" stroke={1.5} />,
    pathname: '/dashboard/products',
  },
]

export function useActiveLink() {
  const router = useRouter()
  return links.find((link) => link.pathname === router.pathname)
}
