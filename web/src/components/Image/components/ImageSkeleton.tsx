import { Skeleton } from '@mantine/core'
import { ImageProps } from 'next/image'
import React from 'react'

type ImageSkeletonProps = Pick<ImageProps, 'width' | 'height'> & {
  style?: React.CSSProperties
}

export function ImageSkeleton({ height, width }: ImageSkeletonProps) {
  return (
    <Skeleton>
      <svg width={width} height={height} viewBox={`0 0 ${width} ${height}`} />
    </Skeleton>
  )
}
