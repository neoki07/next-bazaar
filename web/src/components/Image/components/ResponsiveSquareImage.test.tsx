import { render, screen } from '@testing-library/react'
import { ResponsiveSquareImage } from './ResponsiveSquareImage'

describe('ResponsiveSquareImage', () => {
  it('renders an image with the correct src and alt attributes', () => {
    const src = 'https://example.com/image.jpg'
    const alt = 'Example Image'
    render(<ResponsiveSquareImage src={src} alt={alt} />)

    const image = screen.getByAltText(alt)
    expect(image).toBeInTheDocument()
    expect(image.getAttribute('src')).toContain(
      `/_next/image?url=${encodeURIComponent(src)}`
    )
    expect(image).toHaveAttribute('alt', alt)
  })

  it('renders a skeleton when isLoading is true and the image is still loading', () => {
    const src = 'https://example.com/image.jpg'
    const alt = 'Example Image'
    render(<ResponsiveSquareImage src={src} alt={alt} isLoading />)

    expect(screen.getByLabelText('Image loading skeleton')).toBeInTheDocument()
  })

  it('does not render a skeleton when isLoading is false', () => {
    const src = 'https://example.com/image.jpg'
    const alt = 'Example Image'
    render(<ResponsiveSquareImage src={src} alt={alt} isLoading={false} />)

    expect(
      screen.queryByLabelText('Image loading skeleton')
    ).not.toBeInTheDocument()
  })
})
