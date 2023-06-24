import { MainLayout } from '@/components/Layout'
import { CartProductList, useCart } from '@/features/cart'
import { Container, Group, Title, rem } from '@mantine/core'

export default function CartPage() {
  const { data: cart, isLoading } = useCart()

  return (
    <MainLayout>
      <Container size={1200} miw={800}>
        <Title mb="lg" order={1}>
          Cart
        </Title>
        <Group align="start" spacing={rem(48)}>
          <div style={{ flex: '1' }}>
            <CartProductList
              cartProducts={cart?.products}
              isLoading={isLoading}
            />
          </div>
        </Group>
      </Container>
    </MainLayout>
  )
}
