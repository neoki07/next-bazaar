import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import Decimal from 'decimal.js'
import { Category } from '../types'
import { ProductForm } from './ProductForm'

const user = userEvent.setup()

const allCategories: Category[] = [
  { id: '1', name: 'Category 1' },
  { id: '2', name: 'Category 2' },
]

describe('ProductForm', () => {
  it('renders all form fields', () => {
    render(<ProductForm allCategories={allCategories} onSubmit={jest.fn()} />)

    expect(screen.getByLabelText('Name')).toBeInTheDocument()
    expect(screen.getByLabelText('Description')).toBeInTheDocument()
    expect(screen.getByLabelText('Category')).toBeInTheDocument()
    expect(screen.getByLabelText('Price')).toBeInTheDocument()
    expect(screen.getByLabelText('StockQuantity')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument()
  })

  it('renders all categories as options', () => {
    render(<ProductForm allCategories={allCategories} onSubmit={jest.fn()} />)

    const categorySelect = screen.getByLabelText('Category')
    expect(categorySelect).toHaveDisplayValue('')
    expect(within(categorySelect).getByRole('option', { name: 'Category 1' }))
    expect(within(categorySelect).getByRole('option', { name: 'Category 2' }))
  })

  it('shows validation errors when submitting an empty form', async () => {
    const handleSubmit = jest.fn()
    render(
      <ProductForm allCategories={allCategories} onSubmit={handleSubmit} />
    )

    await user.click(screen.getByRole('button', { name: 'Save' }))

    expect(await screen.findAllByText('Required')).toHaveLength(4)
    expect(handleSubmit).not.toHaveBeenCalled()
  })

  it('submits the form with valid data', async () => {
    const handleSubmit = jest.fn()
    render(
      <ProductForm allCategories={allCategories} onSubmit={handleSubmit} />
    )

    await user.type(screen.getByLabelText('Name'), 'Product 1')
    await user.type(screen.getByLabelText('Description'), 'Description 1')
    await user.selectOptions(screen.getByLabelText('Category'), '1')
    await user.type(screen.getByLabelText('Price'), '10.00')
    await user.type(screen.getByLabelText('StockQuantity'), '5')
    await user.click(screen.getByRole('button', { name: 'Save' }))

    expect(screen.queryByText('Required')).not.toBeInTheDocument()
    expect(handleSubmit).toHaveBeenCalledWith({
      name: 'Product 1',
      description: 'Description 1',
      categoryId: '1',
      price: new Decimal('10.00'),
      stockQuantity: 5,
    })
  })
})
