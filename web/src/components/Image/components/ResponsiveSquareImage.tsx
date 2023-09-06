import { AspectRatio, Skeleton } from '@mantine/core'
import NextImage, { ImageProps as NextImageProps } from 'next/image'
import { ReactNode, useCallback, useState } from 'react'

interface LoadingImageProps extends Omit<NextImageProps, 'src' | 'alt'> {
  src?: string
  alt?: string
  isLoading: true
}

interface NonLoadingImageProps extends Omit<NextImageProps, 'src' | 'alt'> {
  src: string
  alt: string
  isLoading?: false
}

type ImageProps = LoadingImageProps | NonLoadingImageProps

export function ResponsiveSquareImage({
  src,
  alt,
  fill = true,
  onLoadingComplete,
  isLoading,
  ...props
}: ImageProps) {
  const [isImageLoading, setIsImageLoading] = useState(true)

  const handleLoadingComplete = useCallback(
    (img: HTMLImageElement) => {
      setIsImageLoading(false)
      onLoadingComplete?.(img)
    },
    [onLoadingComplete]
  )

  return (
    <Wrapper isLoading={isLoading || isImageLoading}>
      <AspectRatio ratio={1}>
        {!isLoading && (
          <NextImage
            {...props}
            src={src}
            alt={alt}
            fill={fill}
            onLoadingComplete={handleLoadingComplete}
          />
        )}
      </AspectRatio>
    </Wrapper>
  )
}

function Wrapper({
  children,
  isLoading,
}: {
  children: ReactNode
  isLoading: boolean
}) {
  if (isLoading) {
    return <Skeleton aria-label="Image loading skeleton">{children}</Skeleton>
  }

  return <>{children}</>
}
