import { zodResolver } from '@hookform/resolvers/zod'
import { BoxProps } from '@mantine/core'
import { useId } from '@mantine/hooks'
import { ReactNode, useState } from 'react'
import {
  FieldValues,
  FormProvider,
  SubmitErrorHandler,
  SubmitHandler,
  UseFormProps as UseHookFormProps,
  useForm as useHookForm,
} from 'react-hook-form'
import { z } from 'zod'

type AsyncDefaultValues<TFieldValues> = (
  payload?: unknown
) => Promise<TFieldValues>

type UseFormProps<TFieldValues extends FieldValues, TContext> = Omit<
  UseHookFormProps<TFieldValues, TContext>,
  'defaultValues'
> & {
  defaultValues: TFieldValues | AsyncDefaultValues<TFieldValues>
  schema?: z.ZodType<TFieldValues>
  onSubmit: SubmitHandler<TFieldValues>
  onSubmitError?: SubmitErrorHandler<TFieldValues>
}

type FormProps = {
  children?: ReactNode
} & Omit<BoxProps, 'children'>

export function useForm<
  TFieldValues extends FieldValues = FieldValues,
  TContext = any
>(props: UseFormProps<TFieldValues, TContext>) {
  const id = useId()

  const { schema, defaultValues, onSubmit, onSubmitError, ...rest } = props

  const methods = useHookForm<TFieldValues, TContext>({
    resolver: schema ? zodResolver(schema) : undefined,
    defaultValues: defaultValues as UseHookFormProps<
      TFieldValues,
      TContext
    >['defaultValues'],
    ...rest,
  })

  const [Form] = useState(() => {
    function Form({ children }: FormProps) {
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

    return Form
  })

  return [Form, methods] as const
}
