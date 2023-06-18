import { MainLayout } from '@/components/Layout'
import { CartProductList, useCartProducts } from '@/features/cart'
import { Container, Group, Title, rem } from '@mantine/core'

export default function CartPage() {
  const { data: cartProducts } = useCartProducts()

  return (
    <MainLayout>
      <Container size={1200} miw={800}>
        <Title mb="lg" order={1}>
          Cart
        </Title>
        <Group align="start" spacing={rem(48)}>
          <div style={{ flex: '1' }}>
            {cartProducts && <CartProductList cartProducts={cartProducts} />}
          </div>
        </Group>
      </Container>
    </MainLayout>
  )
}
