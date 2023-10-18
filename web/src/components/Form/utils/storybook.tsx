import { useForm } from '@/components/Form'
import { PartialStoryFn } from '@storybook/csf'
import { ReactRenderer } from '@storybook/react'
import { ReactNode, useEffect, useRef } from 'react'
import { Resolver } from 'react-hook-form'

// TODO: Optimize the types for `resolver` and `defaultValues`.

function Wrapper({
  children,
  resolver,
  defaultValues,
  submitOnMount,
}: {
  children: ReactNode
  resolver: Resolver
  defaultValues: any
  submitOnMount: boolean
}) {
  const submitOnMountRef = useRef(submitOnMount)
  const [Form, methods] = useForm({
    resolver,
    defaultValues,
  })

  useEffect(() => {
    if (submitOnMountRef.current) {
      methods.trigger()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return <Form>{children}</Form>
}

export function renderDecorator(
  Story: PartialStoryFn<ReactRenderer>,
  resolver: Resolver,
  defaultValues: any,
  submitOnMount: boolean = false
) {
  return (
    <Wrapper
      resolver={resolver}
      defaultValues={defaultValues}
      submitOnMount={submitOnMount}
    >
      <Story />
    </Wrapper>
  )
}
