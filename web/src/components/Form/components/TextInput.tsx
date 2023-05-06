import { TextInput as MantineTextInput } from '@mantine/core'
import { TextInputProps } from '../types'
import { useController } from 'react-hook-form'
import { ErrorMessage } from './ErrorMessage'
import { FC } from 'react'

export const TextInput: FC<TextInputProps> = ({ label, name, ...rest }) => {
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
