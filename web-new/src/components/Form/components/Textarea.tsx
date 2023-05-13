import { Textarea as MantineTextarea } from '@mantine/core'
import { useController } from 'react-hook-form'
import { TextareaProps } from '../types'
import { ErrorMessage } from './ErrorMessage'

export function Textarea({ label, name, ...rest }: TextareaProps) {
  const {
    field,
    fieldState: { error: fieldError },
  } = useController({ name })

  const error = fieldError ? (
    <ErrorMessage>{fieldError.message?.toString()}</ErrorMessage>
  ) : undefined

  return (
    <MantineTextarea
      id={name}
      label={label}
      error={error}
      {...rest}
      {...field}
    />
  )
}
