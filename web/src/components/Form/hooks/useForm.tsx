import {
  Button,
  ButtonProps,
  useMantineTheme,
  GridProps,
  BoxProps,
} from '@mantine/core'
import { useId } from '@mantine/hooks'
import React, { FC, ReactNode, useMemo } from 'react'
import {
  useForm as useHookForm,
  FormProvider,
  FieldValues,
  UseFormProps,
  SubmitHandler,
  SubmitErrorHandler,
} from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'

type AsyncDefaultValues<TFieldValues> = (
  payload?: unknown
) => Promise<TFieldValues>

type FormProps<TFieldValues extends FieldValues, TContext> = Omit<
  UseFormProps<TFieldValues, TContext>,
  'defaultValues'
> & {
  defaultValues: TFieldValues | AsyncDefaultValues<TFieldValues>
  schema?: z.ZodType<TFieldValues>
  onSubmit: SubmitHandler<TFieldValues>
  onSubmitError?: SubmitErrorHandler<TFieldValues>
}

export const useForm = <
  TFieldValues extends FieldValues = FieldValues,
  TContext = any
>(
  props: FormProps<TFieldValues, TContext>
) => {
  const id = useId()
  const theme = useMantineTheme()

  const { schema, defaultValues, onSubmit, onSubmitError, ...rest } = props

  const methods = useHookForm<TFieldValues, TContext>({
    resolver: schema ? zodResolver(schema) : undefined,
    defaultValues: defaultValues as UseFormProps<
      TFieldValues,
      TContext
    >['defaultValues'],
    ...rest,
  })

  const Form = useMemo(() => {
    const Form: FC<
      {
        children?: ReactNode
        grid?: Omit<GridProps, 'children'>
      } & Omit<BoxProps, 'children'>
    > & { SubmitButton: FC<ButtonProps> } = ({ children, grid, ...rest }) => {
      return (
        <FormProvider {...methods}>
          <form
            id={id}
            onSubmit={methods.handleSubmit(onSubmit, onSubmitError)}
          >
            {children}
          </form>
        </FormProvider>
      )
    }

    const SubmitButton: FC<ButtonProps> = (props) => (
      <Button
        type="submit"
        form={id}
        loaderProps={{ color: theme.colors.blue[5] }}
        {...props}
      />
    )

    Form.SubmitButton = SubmitButton

    return Form
  }, [])

  return [Form, methods] as const
}
