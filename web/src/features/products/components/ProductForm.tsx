import {
  NativeSelect,
  NumberInput,
  TextInput,
  Textarea,
  useForm,
} from '@/components/Form'
import { uploadFile } from '@/features/upload'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button, Center, Grid, Stack } from '@mantine/core'
import Decimal from 'decimal.js'
import { useState } from 'react'
import { SubmitHandler } from 'react-hook-form'
import { z } from 'zod'
import { Category, Product } from '../types'
import { ProductImageDropzone } from './ProductImageDropzone'

export interface ProductFormValues {
  name: string
  description?: string
  categoryId: string
  price: Decimal
  stockQuantity: number
  imageUrl?: string
}

export function productToFormValues(product: Product): ProductFormValues {
  return {
    name: product.name,
    description: product.description,
    categoryId: product.categoryId,
    price: product.price,
    stockQuantity: product.stockQuantity,
    imageUrl: product.imageUrl,
  }
}

interface ProductFormProps {
  imageSize: number
  allCategories: Category[]
  initialValues?: ProductFormValues
  onSubmit: SubmitHandler<ProductFormValues>
}

export function ProductForm({
  imageSize,
  allCategories,
  initialValues,
  onSubmit,
}: ProductFormProps) {
  const [uploadingImage, setUploadingImage] = useState(false)

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
    imageUrl: z.string().optional(),
  })

  const [Form, methods] = useForm<{
    name: string
    description?: string
    categoryId: string
    price?: number
    stockQuantity?: number
    imageUrl?: string
  }>({
    resolver: zodResolver(schema),
    defaultValues:
      initialValues === undefined
        ? {
            name: '',
            description: undefined,
            categoryId: '',
            price: undefined,
            stockQuantity: undefined,
            imageUrl: undefined,
          }
        : {
            ...initialValues,
            price: initialValues.price.toNumber(),
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

        <ProductImageDropzone
          label="Image"
          name="imageUrl"
          loading={uploadingImage}
          imageWidth={imageSize}
          imageHeight={imageSize}
          uploadedImageUrl={methods.watch('imageUrl')}
          onDrop={(file) => {
            setUploadingImage(true)
            uploadFile(file)
              .then((url) => {
                methods.setValue('imageUrl', url)
                setUploadingImage(false)
              })
              .catch((error) => {
                throw new Error(error)
              })
          }}
          onRemove={() => {
            methods.setValue('imageUrl', undefined)
          }}
        />

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

        <Grid>
          <Grid.Col span={6}>
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
          </Grid.Col>

          <Grid.Col span={6}>
            <NumberInput
              label="StockQuantity"
              name="stockQuantity"
              withAsterisk
              min={0}
            />
          </Grid.Col>
        </Grid>

        <Center mt="sm">
          <Button
            type="submit"
            loading={methods.formState.isSubmitting}
            disabled={uploadingImage}
          >
            Save
          </Button>
        </Center>
      </Stack>
    </Form>
  )
}
