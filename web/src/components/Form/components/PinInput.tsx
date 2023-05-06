import { PinInputProps } from '../types'
import { Input, PinInput as MantinePinInput } from '@mantine/core'
import { useController } from 'react-hook-form'
import { ErrorMessage } from './ErrorMessage'
import { FC } from 'react'

export const PinInput: FC<PinInputProps> = ({
  label,
  name,
  description,
  descriptionProps,
  required,
  withAsterisk,
  labelProps,
  errorProps,
  inputContainer,
  inputWrapperOrder,
  ...rest
}) => {
  const {
    field,
    fieldState: { error: fieldError },
  } = useController({ name })

  const errorMessage = fieldError ? (
    <ErrorMessage>{fieldError.message?.toString()}</ErrorMessage>
  ) : undefined

  return (
    <Input.Wrapper
      id={name}
      label={label}
      error={errorMessage}
      {...{
        description,
        descriptionProps,
        required,
        withAsterisk,
        labelProps,
        errorProps,
        inputContainer,
        inputWrapperOrder,
      }}
    >
      <MantinePinInput id={name} error={!!errorMessage} {...rest} {...field} />
    </Input.Wrapper>
  )
}
