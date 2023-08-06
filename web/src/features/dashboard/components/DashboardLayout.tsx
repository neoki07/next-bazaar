import { MainLayout } from '@/components/Layout'
import { Container, Flex, rem } from '@mantine/core'
import { ReactNode } from 'react'
import { NavBar } from './NavBar'

interface DashboardLayoutProps {
  children: ReactNode
}

export function DashboardLayout({ children }: DashboardLayoutProps) {
  return (
    <MainLayout>
      <Container>
        <Flex gap={rem(40)}>
          <NavBar width={200} />
          <div style={{ flex: 1 }}>{children}</div>
        </Flex>
      </Container>
    </MainLayout>
  )
}
