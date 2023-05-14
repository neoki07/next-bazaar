import { formatMoneyFromDecimal } from '@/utils/money'
import { Text } from '@mantine/core'
import Decimal from 'decimal.js'

interface PriceProps {
  price: Decimal
}

export function Price({ price }: PriceProps) {
  return (
    <Text size="xl" weight="bold">
      {formatMoneyFromDecimal(price)}
    </Text>
  )
}
