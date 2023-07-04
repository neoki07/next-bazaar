import { formatMoneyFromDecimal } from '@/utils/money'
import { Text } from '@mantine/core'
import Decimal from 'decimal.js'

interface PriceProps {
  price: Decimal
  className?: string
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | string | number
  weight?: React.CSSProperties['fontWeight']
}

export function Price({ price, className, size, weight }: PriceProps) {
  return (
    <Text className={className} size={size} weight={weight}>
      {formatMoneyFromDecimal(price)}
    </Text>
  )
}
