import {
  NativeSelect,
  NumberInput,
  TextInput,
  Textarea,
  useForm,
} from '@/components/Form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button, Center, Stack } from '@mantine/core'
import Decimal from 'decimal.js'
import { SubmitHandler } from 'react-hook-form'
import { z } from 'zod'
import { Category } from '../types'

export interface ProductFormValues {
  name: string
  description?: string
  categoryId: string
  price: Decimal
  stockQuantity: number
}

interface ProductFormProps {
  allCategories: Category[]
  onSubmit: SubmitHandler<ProductFormValues>
}

export function ProductForm({ allCategories, onSubmit }: ProductFormProps) {
  const schema = z.object({
    name: z.string().min(1, { message: 'Required' }),
    description: z
      .string()
      .optional()
      .transform((value) => (value === '' ? undefined : value)),
    categoryId: z
      .string()
      .min(1, { message: 'Required' })
      .refine((value) =>
        allCategories.some((category) => category.id === value)
      ),
    price: z.number({ required_error: 'Required' }).min(0.01),
    stockQuantity: z.number({ required_error: 'Required' }).min(0),
  })

  const [Form, methods] = useForm<{
    name: string
    description?: string
    categoryId: string
    price?: number
    stockQuantity?: number
  }>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: '',
      description: undefined,

      categoryId: '',
      price: undefined,
      stockQuantity: undefined,
    },
    onSubmit: (data) => {
      if (data.price === undefined) {
        throw new Error('Price is undefined')
      }

      if (data.stockQuantity === undefined) {
        throw new Error('StockQuantity is undefined')
      }

      onSubmit({
        ...data,
        price: new Decimal(data.price),
        stockQuantity: data.stockQuantity,
      })
    },
  })

  return (
    <Form>
      <Stack>
        <TextInput label="Name" name="name" withAsterisk />

        <Textarea label="Description" name="description" minRows={5} />

        <NativeSelect
          label="Category"
          name="categoryId"
          withAsterisk
          options={[
            { label: '', value: '' },
            ...allCategories.map((category) => ({
              label: category.name,
              value: category.id,
            })),
          ]}
        />

        <NumberInput
          label="Price"
          name="price"
          withAsterisk
          precision={2}
          min={0.01}
          step={0.01}
          parser={(value) => value.replace(/\$\s?|(,*)/g, '')}
          formatter={(value) =>
            !Number.isNaN(parseFloat(value))
              ? `$ ${value}`.replace(/\B(?<!\.\d*)(?=(\d{3})+(?!\d))/g, ',')
              : '$ '
          }
        />

        <NumberInput
          label="StockQuantity"
          name="stockQuantity"
          withAsterisk
          min={0}
        />

        <Center mt="sm">
          <Button type="submit" loading={methods.formState.isSubmitting}>
            Save
          </Button>
        </Center>
      </Stack>
    </Form>
  )
}
