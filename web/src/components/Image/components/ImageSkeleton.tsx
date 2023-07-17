import { Skeleton } from '@mantine/core'
import { ImageProps } from 'next/image'
import React from 'react'

interface ImageSkeletonProps extends Pick<ImageProps, 'width' | 'height'> {
  style?: React.CSSProperties
}

export function ImageSkeleton({ width, height }: ImageSkeletonProps) {
  return (
    <Skeleton style={{ width, height }}>
      <svg width={width} height={height} viewBox={`0 0 ${width} ${height}`} />
    </Skeleton>
  )
}
