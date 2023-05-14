import { Input, PinInput as MantinePinInput } from '@mantine/core'
import { useController } from 'react-hook-form'
import { PinInputProps } from '../types'
import { ErrorMessage } from './ErrorMessage'

export function PinInput({
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
}: PinInputProps) {
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
