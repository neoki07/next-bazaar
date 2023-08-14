import { render, screen, within } from '@testing-library/react'
import Decimal from 'decimal.js'
import { ProductList } from './ProductList'

const IMAGE_SIZE = 200

const products = [
  {
    id: '1',
    name: 'Product 1',
    description: 'Description 1',
    price: new Decimal(10.0),
    stockQuantity: 10,
    category: 'Category 1',
    seller: 'Seller 1',
    imageUrl: 'https://example.com/image.png',
  },
  {
    id: '2',
    name: 'Product 2',
    description: 'Description 2',
    price: new Decimal(20.0),
    stockQuantity: 10,
    category: 'Category 2',
    seller: 'Seller 2',
    imageUrl: 'https://example.com/image.png',
  },
]

describe('ProductList', () => {
  it('renders products', () => {
    render(<ProductList products={products} imageSize={IMAGE_SIZE} />)

    const list = screen.getByRole('list')
    expect(list).toBeInTheDocument()
    expect(within(list).getAllByRole('listitem')).toHaveLength(2)
  })

  it('renders "No products" when products is empty', () => {
    render(<ProductList products={[]} imageSize={IMAGE_SIZE} />)

    expect(screen.getByText('No products')).toBeInTheDocument()
  })

  it('renders skeleton when isLoading is true', () => {
    render(<ProductList imageSize={IMAGE_SIZE} isLoading />)

    const list = screen.getByRole('list')
    expect(list).toBeInTheDocument()
    expect(within(list).getAllByRole('listitem')).toHaveLength(3)
  })

  it('renders skeleton when isLoading is true even when products is not empty', () => {
    render(<ProductList products={products} imageSize={IMAGE_SIZE} isLoading />)

    const list = screen.getByRole('list')
    expect(list).toBeInTheDocument()
    expect(within(list).getAllByRole('listitem')).toHaveLength(3)
    expect(within(list).queryByText('Product 1')).not.toBeInTheDocument()
    expect(within(list).queryByText('Product 2')).not.toBeInTheDocument()
  })
})
