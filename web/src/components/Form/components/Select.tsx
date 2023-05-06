import { Select as MantineSelect } from '@mantine/core'
import { IconChevronDown } from '@tabler/icons-react'
import { SelectProps } from '../types'
import { useController } from 'react-hook-form'
import { ErrorMessage } from './ErrorMessage'
import { FC } from 'react'

export const Select: FC<SelectProps> = ({ label, options, name, ...rest }) => {
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
    <MantineSelect
      id={name}
      rightSection={<IconChevronDown width={15} color="#9e9e9e" />}
      styles={{ rightSection: { pointerEvents: 'none' } }}
      label={label}
      onChange={(value) => onChange(value ?? defaultValues?.[name])}
      allowDeselect
      error={error}
      {...rest}
      data={options}
      {...restField}
    />
  )
}
