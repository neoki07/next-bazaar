import Decimal from 'decimal.js'
import { Cart } from '../types'
import { transformCart } from './transform'

describe('transformCart', () => {
  it('transforms the cart correctly', () => {
    const cart = {
      products: [
        {
          id: '1',
          name: 'Product 1',
          description: 'Description 1',
          price: '10.00',
          quantity: 2,
          subtotal: '20.00',
          image_url: 'https://via.placeholder.com/200',
        },
        {
          id: '2',
          name: 'Product 2',
          description: 'Description 2',
          price: '5.00',
          quantity: 1,
          subtotal: '5.00',
          image_url: 'https://via.placeholder.com/200',
        },
      ],
      subtotal: '25.00',
      shipping: '5.00',
      tax: '2.50',
      total: '32.50',
    }

    const expectedCart: Cart = {
      products: [
        {
          id: '1',
          name: 'Product 1',
          description: 'Description 1',
          price: new Decimal('10.00'),
          quantity: 2,
          subtotal: new Decimal('20.00'),
          imageUrl: 'https://via.placeholder.com/200',
        },
        {
          id: '2',
          name: 'Product 2',
          description: 'Description 2',
          price: new Decimal('5.00'),
          quantity: 1,
          subtotal: new Decimal('5.00'),
          imageUrl: 'https://via.placeholder.com/200',
        },
      ],
      subtotal: new Decimal('25.00'),
      shipping: new Decimal('5.00'),
      tax: new Decimal('2.50'),
      total: new Decimal('32.50'),
    }

    const transformedCart = transformCart(cart)

    expect(transformedCart).toEqual(expectedCart)
  })

  it('transforms the cart correctly if optional fields are undefined', () => {
    const cart = {
      products: [
        {
          id: '1',
          name: 'Product 1',
          price: '10.00',
          quantity: 2,
          subtotal: '20.00',
        },
      ],
      subtotal: '10.00',
      shipping: '5.00',
      tax: '1.00',
      total: '16.00',
    }

    const expectedCart: Cart = {
      products: [
        {
          id: '1',
          name: 'Product 1',
          price: new Decimal('10.00'),
          quantity: 2,
          subtotal: new Decimal('20.00'),
        },
      ],
      subtotal: new Decimal('10.00'),
      shipping: new Decimal('5.00'),
      tax: new Decimal('1.00'),
      total: new Decimal('16.00'),
    }

    const transformedCart = transformCart(cart)

    expect(transformedCart).toEqual(expectedCart)
  })

  it('throws an error if product id is undefined', () => {
    const cart = {
      products: [
        {
          name: 'Product 1',
          price: '10.00',
          quantity: 2,
          subtotal: '20.00',
        },
      ],
      subtotal: '10.00',
      shipping: '5.00',
      tax: '1.00',
      total: '16.00',
    }

    expect(() => transformCart(cart)).toThrowError(
      'required fields are undefined'
    )
  })

  it('throws an error if product name is undefined', () => {
    const cart = {
      products: [
        {
          id: '1',
          price: '10.00',
          quantity: 2,
          subtotal: '20.00',
        },
      ],
      subtotal: '10.00',
      shipping: '5.00',
      tax: '1.00',
      total: '16.00',
    }

    expect(() => transformCart(cart)).toThrowError(
      'required fields are undefined'
    )
  })

  it('throws an error if product price is undefined', () => {
    const cart = {
      products: [
        {
          id: '1',
          name: 'Product 1',
          quantity: 2,
          subtotal: '20.00',
        },
      ],
      subtotal: '10.00',
      shipping: '5.00',
      tax: '1.00',
      total: '16.00',
    }

    expect(() => transformCart(cart)).toThrowError(
      'required fields are undefined'
    )
  })

  it('throws an error if product quantity is undefined', () => {
    const cart = {
      products: [
        {
          id: '1',
          name: 'Product 1',
          price: '10.00',
          subtotal: '20.00',
        },
      ],
      subtotal: '10.00',
      shipping: '5.00',
      tax: '1.00',
      total: '16.00',
    }

    expect(() => transformCart(cart)).toThrowError(
      'required fields are undefined'
    )
  })

  it('throws an error if product subtotal undefined', () => {
    const cart = {
      products: [
        {
          id: '1',
          name: 'Product 1',
          price: '10.00',
          quantity: 2,
        },
      ],
      subtotal: '10.00',
      shipping: '5.00',
      tax: '1.00',
      total: '16.00',
    }

    expect(() => transformCart(cart)).toThrowError(
      'required fields are undefined'
    )
  })

  it('throws an error if subtotal is undefined', () => {
    const cart = {
      products: [
        {
          id: '1',
          name: 'Product 1',
          price: '10.00',
          quantity: 2,
          subtotal: '20.00',
        },
      ],
      shipping: '5.00',
      tax: '1.00',
      total: '16.00',
    }

    expect(() => transformCart(cart)).toThrowError(
      'required fields are undefined'
    )
  })

  it('throws an error if shipping is undefined', () => {
    const cart = {
      products: [
        {
          id: '1',
          name: 'Product 1',
          price: '10.00',
          quantity: 2,
          subtotal: '20.00',
        },
      ],
      subtotal: '10.00',
      tax: '1.00',
      total: '16.00',
    }

    expect(() => transformCart(cart)).toThrowError(
      'required fields are undefined'
    )
  })

  it('throws an error if tax is undefined', () => {
    const cart = {
      products: [
        {
          id: '1',
          name: 'Product 1',
          price: '10.00',
          quantity: 2,
          subtotal: '20.00',
        },
      ],
      subtotal: '10.00',
      shipping: '5.00',
      total: '16.00',
    }

    expect(() => transformCart(cart)).toThrowError(
      'required fields are undefined'
    )
  })

  it('throws an error if total is undefined', () => {
    const cart = {
      products: [
        {
          id: '1',
          name: 'Product 1',
          price: '10.00',
          quantity: 2,
          subtotal: '20.00',
        },
      ],
      subtotal: '10.00',
      shipping: '5.00',
      tax: '1.00',
    }

    expect(() => transformCart(cart)).toThrowError(
      'required fields are undefined'
    )
  })
})
