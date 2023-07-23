import { render, screen } from '@testing-library/react'
import Decimal from 'decimal.js'
import { Product, ProductCard } from '../../products'

const product: Product = {
  id: '1',
  name: 'Product',
  category: 'Category',
  price: new Decimal(10.0),
  stockQuantity: 5,
  seller: 'Seller',
  imageUrl: 'https://example.com/image.png',
}

const getProductLink = (product: Product) => `/products/${product.id}`
const imageSize = 260

describe('ProductCard', () => {
  it('renders product information', () => {
    render(
      <ProductCard
        product={product}
        getProductLink={getProductLink}
        imageSize={imageSize}
      />
    )

    expect(screen.getByText('Product')).toBeInTheDocument()
    expect(screen.getByText('Category')).toBeInTheDocument()
    expect(screen.getByText('$10.00')).toBeInTheDocument()
    expect(screen.getByText('Seller')).toBeInTheDocument()
    const imageElement = screen.getByRole('img', { name: 'Product' })
    expect(imageElement).toBeInTheDocument()
    expect(imageElement.getAttribute('src')).toContain(
      `/_next/image?url=${encodeURIComponent('https://example.com/image.png')}`
    )
  })
})
