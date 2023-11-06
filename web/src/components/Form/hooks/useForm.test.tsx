import { render, renderHook, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { useForm } from './useForm'

const user = userEvent.setup()

describe('useForm', () => {
  it('calls onSubmit when the form is submitted', async () => {
    const onSubmit = jest.fn()
    const { result } = renderHook(() =>
      useForm({
        defaultValues: {},
        onSubmit,
      })
    )
    const [Form] = result.current

    render(
      <Form>
        <button type="submit">Submit</button>
      </Form>
    )

    const button = screen.getByRole('button', { name: 'Submit' })
    await user.click(button)

    expect(onSubmit).toHaveBeenCalled()
  })
})
