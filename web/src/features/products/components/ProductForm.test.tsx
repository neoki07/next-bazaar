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

const imageSize = 100
const noop = () => {}

describe('ProductForm', () => {
  it('renders all form fields', () => {
    render(
      <ProductForm
        imageSize={imageSize}
        allCategories={allCategories}
        onSubmit={jest.fn()}
        onCancel={noop}
      />
    )

    expect(
      screen.getByLabelText((content) => content.includes('Name'))
    ).toBeInTheDocument()
    expect(screen.getByLabelText('Description')).toBeInTheDocument()
    expect(
      screen.getByLabelText((content) => content.includes('Category'))
    ).toBeInTheDocument()
    expect(
      screen.getByLabelText((content) => content.includes('Price'))
    ).toBeInTheDocument()
    expect(
      screen.getByLabelText((content) => content.includes('StockQuantity'))
    ).toBeInTheDocument()

    expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument()
  })

  it('renders all categories as options', () => {
    render(
      <ProductForm
        imageSize={imageSize}
        allCategories={allCategories}
        onSubmit={jest.fn()}
        onCancel={noop}
      />
    )

    const categorySelect = screen.getByLabelText((content) =>
      content.includes('Category')
    )
    expect(categorySelect).toHaveDisplayValue('')
    expect(within(categorySelect).getByRole('option', { name: 'Category 1' }))
    expect(within(categorySelect).getByRole('option', { name: 'Category 2' }))
  })

  it('shows validation errors when submitting an empty form', async () => {
    const handleSubmit = jest.fn()
    render(
      <ProductForm
        imageSize={imageSize}
        allCategories={allCategories}
        onSubmit={handleSubmit}
        onCancel={noop}
      />
    )

    await user.click(screen.getByRole('button', { name: 'Save' }))

    expect(await screen.findAllByText('Required')).toHaveLength(4)
    expect(handleSubmit).not.toHaveBeenCalled()
  })

  it('submits the form with valid data', async () => {
    const handleSubmit = jest.fn()
    render(
      <ProductForm
        imageSize={imageSize}
        allCategories={allCategories}
        onSubmit={handleSubmit}
        onCancel={noop}
      />
    )

    await user.type(
      screen.getByLabelText((content) => content.includes('Name')),
      'Product 1'
    )
    await user.type(screen.getByLabelText('Description'), 'Description 1')
    await user.selectOptions(
      screen.getByLabelText((content) => content.includes('Category')),
      '1'
    )
    await user.type(
      screen.getByLabelText((content) => content.includes('Price')),
      '10.00'
    )
    await user.type(
      screen.getByLabelText((content) => content.includes('StockQuantity')),
      '5'
    )
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

  // TODO: onCancel
})
