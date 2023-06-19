import { NativeSelect as MantineNativeSelect } from '@mantine/core'
import { useController } from 'react-hook-form'
import { NativeSelectProps } from '../types'
import { ErrorMessage } from './ErrorMessage'

export function NativeSelect({
  label,
  options,
  name,
  ...rest
}: NativeSelectProps) {
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
    <MantineNativeSelect
      id={name}
      styles={{ rightSection: { pointerEvents: 'none' } }}
      label={label}
      onChange={(value) => onChange(value ?? defaultValues?.[name])}
      error={error}
      {...rest}
      data={options}
      {...restField}
    />
  )
}
