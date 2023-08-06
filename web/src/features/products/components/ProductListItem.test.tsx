import { render, screen } from '@testing-library/react'
import Decimal from 'decimal.js'
import { ProductListItem } from './ProductListItem'

const IMAGE_SIZE = 200

const product = {
  id: '1',
  name: 'Product',
  description: 'Description',
  price: new Decimal(10.0),
  stockQuantity: 10,
  category: 'Category',
  seller: 'Seller',
  imageUrl: 'https://example.com/image.png',
}

describe('ProductListItem', () => {
  it('renders product information', () => {
    render(<ProductListItem product={product} imageSize={IMAGE_SIZE} />)

    expect(screen.getByText('Product')).toBeInTheDocument()
    expect(screen.getByText('$10.00')).toBeInTheDocument()

    const image = screen.getByRole('img', { name: 'Product' })
    expect(image).toBeInTheDocument()
    expect(image.getAttribute('src')).toContain(
      `/_next/image?url=${encodeURIComponent('https://example.com/image.png')}`
    )
  })
})
