import Decimal from 'decimal.js'

export function formatMoneyFromDecimal(money: Decimal): string {
  const roundedMoney = money.toFixed(2)

  const parts = roundedMoney.split('.')
  const formattedIntegerPart = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, ',')
  const formattedMoney =
    parts.length > 1
      ? formattedIntegerPart + '.' + parts[1]
      : formattedIntegerPart

  return '$' + formattedMoney
}
