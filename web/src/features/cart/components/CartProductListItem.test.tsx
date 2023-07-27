import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import Decimal from 'decimal.js'
import { CartProductListItem } from './CartProductListItem'

const IMAGE_SIZE = 200

const user = userEvent.setup()

const cartProduct = {
  id: '1',
  name: 'Product',
  description: 'Description',
  price: new Decimal(10.0),
  quantity: 5,
  subtotal: new Decimal(50.0),
  imageUrl: 'https://example.com/image.png',
}

describe('CartProductListItem', () => {
  it('renders product information', () => {
    render(
      <CartProductListItem cartProduct={cartProduct} imageSize={IMAGE_SIZE} />
    )

    expect(screen.getByText('Product')).toBeInTheDocument()
    expect(screen.getByText('$10.00')).toBeInTheDocument()
    expect(screen.getByLabelText('Quantity')).toHaveValue('5')

    const image = screen.getByRole('img', { name: 'Product' })
    expect(image).toBeInTheDocument()
    expect(image.getAttribute('src')).toContain(
      `/_next/image?url=${encodeURIComponent('https://example.com/image.png')}`
    )
  })

  it('renders quantity select', () => {
    render(
      <CartProductListItem cartProduct={cartProduct} imageSize={IMAGE_SIZE} />
    )

    const quantitySelect = screen.getByLabelText('Quantity')
    expect(quantitySelect).toBeInTheDocument()
    expect(quantitySelect).toHaveValue('5')
  })

  it('calls onChangeQuantity when quantity changes', async () => {
    const onChangeQuantity = jest.fn()
    render(
      <CartProductListItem
        cartProduct={cartProduct}
        imageSize={IMAGE_SIZE}
        onChangeQuantity={onChangeQuantity}
      />
    )

    await user.selectOptions(screen.getByLabelText('Quantity'), '3')

    expect(onChangeQuantity).toHaveBeenCalledTimes(1)
    expect(onChangeQuantity).toHaveBeenCalledWith('1', 3)
  })

  it('renders delete button', () => {
    render(
      <CartProductListItem cartProduct={cartProduct} imageSize={IMAGE_SIZE} />
    )

    expect(
      screen.getByRole('button', { name: 'Remove product' })
    ).toBeInTheDocument()
  })

  it('calls onDelete when delete button is clicked', async () => {
    const onDelete = jest.fn()
    render(
      <CartProductListItem
        cartProduct={cartProduct}
        imageSize={IMAGE_SIZE}
        onDelete={onDelete}
      />
    )

    await user.click(screen.getByRole('button', { name: 'Remove product' }))

    expect(onDelete).toHaveBeenCalledTimes(1)
    expect(onDelete).toHaveBeenCalledWith('1')
  })

  it('renders quantity text when editable is false', () => {
    render(
      <CartProductListItem
        cartProduct={cartProduct}
        imageSize={IMAGE_SIZE}
        editable={false}
      />
    )

    expect(screen.getByText('Quantity: 5')).toBeInTheDocument()
    expect(screen.queryByLabelText('Quantity')).not.toBeInTheDocument()
  })
})
