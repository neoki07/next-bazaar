import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import Decimal from 'decimal.js'
import { CartProductList } from './CartProductList'

const IMAGE_SIZE = 200

const user = userEvent.setup()

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
      <CartProductList cartProducts={cartProducts} imageSize={IMAGE_SIZE} />
    )

    const list = screen.getByRole('list')
    expect(list).toBeInTheDocument()
    expect(within(list).getAllByRole('listitem')).toHaveLength(2)
  })

  it('renders "No products" when cartProducts is empty', () => {
    render(<CartProductList cartProducts={[]} imageSize={IMAGE_SIZE} />)

    expect(screen.getByText('No products')).toBeInTheDocument()
  })

  it('renders skeleton when isLoading is true', () => {
    render(<CartProductList imageSize={IMAGE_SIZE} isLoading />)

    const list = screen.getByRole('list')
    expect(list).toBeInTheDocument()
    expect(within(list).getAllByRole('listitem')).toHaveLength(3)
  })

  it('renders skeleton when isLoading is true even when cartProducts is not empty', () => {
    render(
      <CartProductList
        cartProducts={cartProducts}
        imageSize={IMAGE_SIZE}
        isLoading
      />
    )

    const list = screen.getByRole('list')
    expect(list).toBeInTheDocument()
    expect(within(list).getAllByRole('listitem')).toHaveLength(3)
    expect(within(list).queryByText('Product 1')).not.toBeInTheDocument()
    expect(within(list).queryByText('Product 2')).not.toBeInTheDocument()
  })

  it('calls onChangeQuantity when quantity changes', async () => {
    const onChangeQuantity = jest.fn()
    render(
      <CartProductList
        cartProducts={cartProducts}
        imageSize={IMAGE_SIZE}
        onChangeQuantity={onChangeQuantity}
      />
    )

    // BUG: 2nd select is not found
    const quantitySelect = screen.getAllByLabelText('Quantity')[0]
    await user.selectOptions(quantitySelect, '3')

    expect(onChangeQuantity).toHaveBeenCalledTimes(1)
    expect(onChangeQuantity).toHaveBeenCalledWith('1', 3)
  })

  it('calls onDelete when delete button is clicked', async () => {
    const onDelete = jest.fn()
    render(
      <CartProductList
        cartProducts={cartProducts}
        imageSize={IMAGE_SIZE}
        onDelete={onDelete}
      />
    )

    const deleteButton = screen.getAllByRole('button', {
      name: 'Remove product',
    })[1]
    await user.click(deleteButton)

    expect(onDelete).toHaveBeenCalledTimes(1)
    expect(onDelete).toHaveBeenCalledWith('2')
  })
})
