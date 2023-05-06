import { Textarea as MantineTextarea } from '@mantine/core'
import { TextareaProps } from '../types'
import { useController } from 'react-hook-form'
import { ErrorMessage } from './ErrorMessage'
import { FC } from 'react'

export const Textarea: FC<TextareaProps> = ({ label, name, ...rest }) => {
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
