import { TextInput as MantineTextInput } from '@mantine/core'
import { useController } from 'react-hook-form'
import { TextInputProps } from '../types'
import { ErrorMessage } from './ErrorMessage'

export function TextInput({ label, name, ...rest }: TextInputProps) {
  const {
    field,
    fieldState: { error: fieldError },
  } = useController({ name })

  const error = fieldError ? (
    <ErrorMessage>{fieldError.message?.toString()}</ErrorMessage>
  ) : undefined

  return (
    <MantineTextInput
      id={name}
      label={label}
      error={error}
      {...rest}
      {...field}
    />
  )
}
