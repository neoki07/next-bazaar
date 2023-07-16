import { render } from '@testing-library/react'
import Decimal from 'decimal.js'
import { Product, ProductCard } from '../../products'

const product: Product = {
  id: '1',
  name: 'test-product',
  category: 'test-category',
  price: new Decimal(10.0),
  stockQuantity: 5,
  seller: 'testuser',
  imageUrl: 'https://example.com/image.png',
}

const getProductLink = (product: Product) => `/products/${product.id}`
const imageSize = 260

describe('ProductCard', () => {
  test('should render product name', () => {
    const { getByText } = render(
      <ProductCard
        product={product}
        getProductLink={getProductLink}
        imageSize={imageSize}
      />
    )

    expect(getByText('test-product')).toBeInTheDocument()
  })
  test('should render category name', () => {
    const { getByText } = render(
      <ProductCard
        product={product}
        getProductLink={getProductLink}
        imageSize={imageSize}
      />
    )

    expect(getByText('test-category')).toBeInTheDocument()
  })
  test('should render price', () => {
    const { getByText } = render(
      <ProductCard
        product={product}
        getProductLink={getProductLink}
        imageSize={imageSize}
      />
    )

    expect(getByText('$10.00')).toBeInTheDocument()
  })
  test('should render seller name', () => {
    const { getByText } = render(
      <ProductCard
        product={product}
        getProductLink={getProductLink}
        imageSize={imageSize}
      />
    )

    expect(getByText('testuser')).toBeInTheDocument()
  })
  test('should render product image', () => {
    const { getByRole } = render(
      <ProductCard
        product={product}
        getProductLink={getProductLink}
        imageSize={imageSize}
      />
    )

    const imageElement = getByRole('img', { name: 'test-product' })
    expect(imageElement).toBeInTheDocument()
    expect(imageElement.getAttribute('src')).toContain(
      `/_next/image?url=${encodeURIComponent('https://example.com/image.png')}`
    )
  })
})
