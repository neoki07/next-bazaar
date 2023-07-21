import { render, screen } from '@testing-library/react'
import { ImageSkeleton } from './ImageSkeleton'

describe('ImageSkeleton', () => {
  it('renders a skeleton with the correct width and height', () => {
    const width = 100
    const height = 200
    render(<ImageSkeleton width={width} height={height} />)

    const skeleton = screen.getByLabelText('Image loading skeleton')
    expect(skeleton).toBeInTheDocument()
    expect(skeleton).toHaveStyle({ width: `${width}px`, height: `${height}px` })
  })
})
