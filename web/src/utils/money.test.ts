import Decimal from 'decimal.js'
import { formatMoneyFromDecimal } from './money'

describe('formatMoneyFromDecimal', () => {
  it('formats a positive decimal with two decimal places', () => {
    const money = new Decimal(1234.56)
    const result = formatMoneyFromDecimal(money)
    expect(result).toEqual('$1,234.56')
  })

  it('formats a negative decimal with two decimal places', () => {
    const money = new Decimal(-1234.56)
    const result = formatMoneyFromDecimal(money)
    expect(result).toEqual('$-1,234.56')
  })

  it('formats a decimal with no decimal places', () => {
    const money = new Decimal(1234)
    const result = formatMoneyFromDecimal(money)
    expect(result).toEqual('$1,234.00')
  })

  it('formats a decimal with more than two decimal places', () => {
    const money = new Decimal(1234.5678)
    const result = formatMoneyFromDecimal(money)
    expect(result).toEqual('$1,234.57')
  })

  it('formats 0 as $0.00', () => {
    const money = new Decimal(0)
    const result = formatMoneyFromDecimal(money)
    expect(result).toEqual('$0.00')
  })
})
