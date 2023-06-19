import { NativeSelect as MantineNativeSelect } from '@mantine/core'
import { useController } from 'react-hook-form'
import { NativeNumberSelectProps } from '../types'
import { ErrorMessage } from './ErrorMessage'

export function NativeNumberSelect({
  label,
  options,
  name,
  ...rest
}: NativeNumberSelectProps) {
  const {
    field,
    fieldState: { error: fieldError },
    formState: { defaultValues },
  } = useController({ name })

  const error = fieldError ? (
    <ErrorMessage>{fieldError.message?.toString()}</ErrorMessage>
  ) : undefined

  const { value, onChange, ...restField } = field

  return (
    <MantineNativeSelect
      id={name}
      styles={{ rightSection: { pointerEvents: 'none' } }}
      label={label}
      value={value === undefined ? '' : value.toString()}
      onChange={(event) => {
        const { value } = event.target
        onChange(
          value === '' ? undefined : Number(value) ?? defaultValues?.[name]
        )
      }}
      error={error}
      {...rest}
      data={options.map((option) => ({
        label: option.toString(),
        value: option.toString(),
      }))}
      {...restField}
    />
  )
}
