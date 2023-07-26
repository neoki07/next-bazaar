import { NativeNumberSelect, useForm } from '@/components/Form'
import { zodResolver } from '@hookform/resolvers/zod'
import { rem } from '@mantine/core'
import { ChangeEvent, useCallback, useEffect } from 'react'
import { z } from 'zod'
import { CartProduct } from '../types'

const schema = z.object({
  quantity: z.number().min(1).max(10),
})

interface QuantitySelectFormProps {
  cartProduct: CartProduct
  options: number[]
  onChange?: (id: string, quantity: number) => void
  disabled?: boolean
}

export function QuantitySelectForm({
  cartProduct,
  options,
  onChange,
  disabled,
}: QuantitySelectFormProps) {
  const [Form, methods] = useForm<{
    quantity: number
  }>({
    resolver: zodResolver(schema),
    defaultValues: {
      quantity: cartProduct.quantity,
    },
  })

  const handleChangeQuantity = useCallback(
    (event: ChangeEvent<HTMLSelectElement>) => {
      const { value } = event.target
      onChange?.(cartProduct.id, Number(value))
    },
    [onChange, cartProduct.id]
  )

  useEffect(() => {
    const { setValue } = methods
    setValue('quantity', cartProduct.quantity)
  }, [methods, cartProduct.quantity])

  return (
    <Form>
      <NativeNumberSelect
        w={rem(80)}
        label="Quantity"
        name="quantity"
        options={options}
        onChange={handleChangeQuantity}
        disabled={disabled}
      />
    </Form>
  )
}
