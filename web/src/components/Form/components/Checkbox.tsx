import { Checkbox as MantineCheckbox } from '@mantine/core'
import { useController } from 'react-hook-form'
import { CheckboxProps } from '../types'
import { ErrorMessage } from './ErrorMessage'

export function Checkbox({ label, name, ...rest }: CheckboxProps) {
  const {
    field,
    fieldState: { error: fieldError },
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
