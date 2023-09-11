import { MainLayout } from '@/components/Layout'
import { useSmallerThan } from '@/hooks'
import { Container, Flex, rem } from '@mantine/core'
import { ReactNode } from 'react'
import { NavBar } from './NavBar'

interface DashboardLayoutProps {
  children: ReactNode
}

export function DashboardLayout({ children }: DashboardLayoutProps) {
  const smallerThanLg = useSmallerThan('sm')

  return (
    <MainLayout>
      <Container>
        {smallerThanLg ? (
          <>{children}</>
        ) : (
          <Flex gap={rem(40)}>
            <NavBar width={200} />
            <div style={{ flex: 1 }}>{children}</div>
          </Flex>
        )}
      </Container>
    </MainLayout>
  )
}
