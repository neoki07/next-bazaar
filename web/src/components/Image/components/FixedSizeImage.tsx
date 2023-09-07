import { Skeleton } from '@mantine/core'
import NextImage, { ImageProps as NativeImageProps } from 'next/image'
import { ReactNode, useCallback, useState } from 'react'

interface LoadingImageProps
  extends Omit<NativeImageProps, 'src' | 'alt' | 'width' | 'height'> {
  src?: string
  alt?: string
  width: number
  height: number
  isLoading: true
}

interface NonLoadingImageProps
  extends Omit<NativeImageProps, 'src' | 'alt' | 'width' | 'height'> {
  src: string
  alt: string
  width: number
  height: number
  isLoading?: false
}

type FixedSizeImageProps = LoadingImageProps | NonLoadingImageProps

export function FixedSizeImage({
  src,
  alt,
  width,
  height,
  onLoadingComplete,
  isLoading,
  ...props
}: FixedSizeImageProps) {
  const [isImageLoading, setIsImageLoading] = useState(true)

  const handleLoadingComplete = useCallback(
    (img: HTMLImageElement) => {
      setIsImageLoading(false)
      onLoadingComplete?.(img)
    },
    [onLoadingComplete]
  )

  return (
    <Wrapper
      isLoading={isLoading || isImageLoading}
      width={width}
      height={height}
    >
      {!isLoading && (
        <NextImage
          {...props}
          src={src}
          alt={alt}
          width={width}
          height={height}
          onLoadingComplete={handleLoadingComplete}
        />
      )}
    </Wrapper>
  )
}

function Wrapper({
  children,
  width,
  height,
  isLoading,
}: {
  children: ReactNode
  width: number
  height: number
  isLoading: boolean
}) {
  if (isLoading) {
    return (
      <Skeleton
        aria-label="Image loading skeleton"
        width={width}
        height={height}
      >
        {children}
      </Skeleton>
    )
  }

  return (
    <div
      style={{
        position: 'relative',
        width,
        height,
      }}
    >
      {children}
    </div>
  )
}
