import { PasswordInput as MantinePasswordInput } from '@mantine/core'
import { PasswordInputProps } from '../types'
import { IconEye, IconEyeOff } from '@tabler/icons-react'
import { useController } from 'react-hook-form'
import { ErrorMessage } from './ErrorMessage'
import { FC } from 'react'

export const PasswordInput: FC<PasswordInputProps> = ({
  label,
  name,
  ...rest
}) => {
  const {
    field,
    fieldState: { error: fieldError },
  } = useController({ name })

  const error = fieldError ? (
    <ErrorMessage>{fieldError.message?.toString()}</ErrorMessage>
  ) : undefined

  return (
    <MantinePasswordInput
      id={name}
      label={label}
      error={error}
      visibilityToggleIcon={({ reveal, size }) =>
        reveal ? <IconEyeOff size={size} /> : <IconEye size={size} />
      }
      {...rest}
      {...field}
    />
  )
}
