import Decimal from 'decimal.js'
import { transformProduct } from './transform'

describe('transformProduct', () => {
  it('transforms the product correctly', () => {
    const product = {
      id: '1',
      name: 'test-product',
      description: 'test-description',
      price: '10.00',
      stock_quantity: 5,
      category: 'test-category',
      seller: 'test-seller',
      image_url: 'https://example.com/image.png',
    }

    const expectedProduct = {
      id: '1',
      name: 'test-product',
      description: 'test-description',
      price: new Decimal('10.00'),
      stockQuantity: 5,
      category: 'test-category',
      seller: 'test-seller',
      imageUrl: 'https://example.com/image.png',
    }

    expect(transformProduct(product)).toEqual(expectedProduct)
  })

  it('transforms the product correctly if optional fields are undefined', () => {
    const product = {
      id: '1',
      name: 'test-product',
      price: '10.00',
      stock_quantity: 5,
      category: 'test-category',
      seller: 'test-seller',
    }

    const expectedProduct = {
      id: '1',
      name: 'test-product',
      price: new Decimal('10.00'),
      stockQuantity: 5,
      category: 'test-category',
      seller: 'test-seller',
    }

    expect(transformProduct(product)).toEqual(expectedProduct)
  })

  it('throws an error if id are undefined', () => {
    const product = {
      name: 'test-product',
      price: '10.00',
      stock_quantity: 5,
      category: 'test-category',
      seller: 'test-seller',
    }

    expect(() => transformProduct(product)).toThrowError(
      'required fields are undefined'
    )
  })

  it('throws an error if name are undefined', () => {
    const product = {
      id: '1',
      price: '10.00',
      stock_quantity: 5,
      category: 'test-category',
      seller: 'test-seller',
    }

    expect(() => transformProduct(product)).toThrowError(
      'required fields are undefined'
    )
  })

  it('throws an error if price are undefined', () => {
    const product = {
      id: '1',
      name: 'test-product',
      stock_quantity: 5,
      category: 'test-category',
      seller: 'test-seller',
    }

    expect(() => transformProduct(product)).toThrowError(
      'required fields are undefined'
    )
  })

  it('throws an error if stock_quantity are undefined', () => {
    const product = {
      id: '1',
      name: 'test-product',
      price: '10.00',
      category: 'test-category',
      seller: 'test-seller',
    }

    expect(() => transformProduct(product)).toThrowError(
      'required fields are undefined'
    )
  })

  it('throws an error if category are undefined', () => {
    const product = {
      id: '1',
      name: 'test-product',
      price: '10.00',
      stock_quantity: 5,
      seller: 'test-seller',
    }

    expect(() => transformProduct(product)).toThrowError(
      'required fields are undefined'
    )
  })

  it('throws an error if seller are undefined', () => {
    const product = {
      id: '1',
      name: 'test-product',
      price: '10.00',
      stock_quantity: 5,
      category: 'test-category',
    }

    expect(() => transformProduct(product)).toThrowError(
      'required fields are undefined'
    )
  })
})
