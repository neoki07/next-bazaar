import { Select as MantineSelect } from '@mantine/core'
import { IconChevronDown } from '@tabler/icons-react'
import { useController } from 'react-hook-form'
import { NumberSelectProps } from '../types'
import { ErrorMessage } from './ErrorMessage'

export function NumberSelect({
  label,
  options,
  name,
  ...rest
}: NumberSelectProps) {
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
    <MantineSelect
      id={name}
      rightSection={<IconChevronDown width={15} color="#9e9e9e" />}
      styles={{ rightSection: { pointerEvents: 'none' } }}
      label={label}
      value={value === undefined ? '' : value.toString()}
      onChange={(value) =>
        onChange(
          value === '' ? undefined : Number(value) ?? defaultValues?.[name]
        )
      }
      allowDeselect
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
