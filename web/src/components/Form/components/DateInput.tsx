import { DateInput as MantineDateInput } from '@mantine/dates'
import { DateInputProps } from '../types'
import { useController } from 'react-hook-form'
import { ErrorMessage } from './ErrorMessage'
import { FC } from 'react'

export const DateInput: FC<DateInputProps> = ({ label, name, ...rest }) => {
  const {
    field,
    fieldState: { error: fieldError },
  } = useController({ name })

  const error = fieldError ? (
    <ErrorMessage>{fieldError.message?.toString()}</ErrorMessage>
  ) : undefined

  return (
    <MantineDateInput
      id={name}
      label={label}
      error={error}
      {...rest}
      {...field}
    />
  )
}
