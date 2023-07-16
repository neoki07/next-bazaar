import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen } from '@testing-library/react'
import Decimal from 'decimal.js'
import { CartProductInfo } from './CartProductInfo'

const queryClient = new QueryClient()

const cartProduct = {
  id: '1',
  name: 'Product',
  description: 'Description',
  price: new Decimal(10.0),
  quantity: 5,
  subtotal: new Decimal(50.0),
  imageUrl: 'https://example.com/image.png',
}

describe('CartProductInfo', () => {
  it('renders product information', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <CartProductInfo cartProduct={cartProduct} imageSize={100} />
      </QueryClientProvider>
    )

    expect(screen.getByText('Product')).toBeInTheDocument()
    expect(screen.getByText('$10.00')).toBeInTheDocument()
    expect(screen.getByLabelText('Quantity')).toHaveValue()
    const imageElement = screen.getByRole('img', { name: 'Product' })
    expect(imageElement).toBeInTheDocument()
    expect(imageElement.getAttribute('src')).toContain(
      `/_next/image?url=${encodeURIComponent('https://example.com/image.png')}`
    )
  })

  // TODO: updates product quantity
  // TODO: deletes product
})
