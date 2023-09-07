import { render, screen } from '@testing-library/react'
import { FixedSizeImage } from './FixedSizeImage'

describe('FixedSizeImage', () => {
  it('renders an image with the correct src and alt attributes', () => {
    const src = 'https://example.com/image.jpg'
    const alt = 'Example Image'
    render(<FixedSizeImage src={src} alt={alt} width={300} height={200} />)

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
    render(
      <FixedSizeImage src={src} alt={alt} width={300} height={200} isLoading />
    )

    const image = screen.getByLabelText('Image loading skeleton')
    expect(image).toBeInTheDocument()
  })
})
