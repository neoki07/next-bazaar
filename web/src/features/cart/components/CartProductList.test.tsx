import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen, within } from '@testing-library/react'
import Decimal from 'decimal.js'
import { CartProductList } from './CartProductList'

const queryClient = new QueryClient()

const cartProducts = [
  {
    id: '1',
    name: 'Product 1',
    price: new Decimal(10.0),
    quantity: 1,
    subtotal: new Decimal(10.0),
  },
  {
    id: '2',
    name: 'Product 2',
    price: new Decimal(20.0),
    quantity: 2,
    subtotal: new Decimal(40.0),
  },
]

describe('CartProductList', () => {
  it('renders cart products', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <CartProductList cartProducts={cartProducts} />
      </QueryClientProvider>
    )

    const list = screen.getByRole('list')
    expect(list).toBeInTheDocument()
    expect(within(list).getAllByRole('listitem')).toHaveLength(2)
  })

  it('renders "No products" when cartProducts is empty', () => {
    render(<CartProductList cartProducts={[]} />)
    expect(screen.getByText('No products')).toBeInTheDocument()
  })

  it('renders skeleton when isLoading is true', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <CartProductList isLoading />
      </QueryClientProvider>
    )

    const list = screen.getByRole('list')
    expect(list).toBeInTheDocument()
    expect(within(list).getAllByRole('listitem')).toHaveLength(3)
  })

  it('renders skeleton when isLoading is true even when cartProducts is not empty', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <CartProductList cartProducts={cartProducts} isLoading />
      </QueryClientProvider>
    )

    const list = screen.getByRole('list')
    expect(list).toBeInTheDocument()
    expect(within(list).getAllByRole('listitem')).toHaveLength(3)
    expect(within(list).queryByText('Product 1')).not.toBeInTheDocument()
    expect(within(list).queryByText('Product 2')).not.toBeInTheDocument()
  })
})
