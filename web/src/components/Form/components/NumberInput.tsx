import { NumberInput as MantineNumberInput } from '@mantine/core'
import { useController } from 'react-hook-form'
import { NumberInputProps } from '../types'
import { ErrorMessage } from './ErrorMessage'

export function NumberInput({ label, name, ...rest }: NumberInputProps) {
  const {
    field,
    fieldState: { error: fieldError },
    formState: { defaultValues },
  } = useController({ name })

  const error = fieldError ? (
    <ErrorMessage>{fieldError.message?.toString()}</ErrorMessage>
  ) : undefined

  const { onChange, ...restField } = field

  return (
    <MantineNumberInput
      id={name}
      label={label}
      onChange={(value) => {
        if (value === '') {
          onChange(defaultValues?.[name] ?? null)
        } else {
          onChange(value)
        }
      }}
      error={error}
      {...rest}
      {...restField}
    />
  )
}
