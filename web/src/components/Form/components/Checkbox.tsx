import { Checkbox as MantineCheckbox } from '@mantine/core'
import { CheckboxProps } from '../types'
import { useController } from 'react-hook-form'
import { ErrorMessage } from './ErrorMessage'
import { FC } from 'react'

export const Checkbox: FC<CheckboxProps> = ({ label, name, ...rest }) => {
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
    <MantineCheckbox
      id={name}
      label={label}
      checked={field.value}
      onChange={(event) => {
        onChange(event.currentTarget.checked)
      }}
      error={error}
      {...rest}
      {...restField}
    />
  )
}
