import { AspectRatio, Skeleton } from '@mantine/core'
import NextImage, { ImageProps as NextImageProps } from 'next/image'

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
  isLoading,
  ...props
}: ImageProps) {
  if (isLoading) {
    return (
      <Skeleton aria-label="Image loading skeleton">
        <AspectRatio ratio={1} />
      </Skeleton>
    )
  }

  return (
    <AspectRatio ratio={1}>
      <NextImage src={src} alt={alt} fill={fill} {...props} />
    </AspectRatio>
  )
}
