import { render, screen } from '@testing-library/react'
import Decimal from 'decimal.js'
import { Price } from './Price'

describe('Price', () => {
  it('renders the price with the correct format', () => {
    const price = new Decimal(12.34)
    render(<Price price={price} />)
    const formattedPrice = screen.getByText('$12.34')
    expect(formattedPrice).toBeInTheDocument()
  })
})
