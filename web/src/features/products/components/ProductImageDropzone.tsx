import { FixedSizeImage } from '@/components/Image'
import {
  Center,
  Group,
  Indicator,
  Input,
  Text,
  UnstyledButton,
  useMantineTheme,
} from '@mantine/core'
import { Dropzone, FileWithPath } from '@mantine/dropzone'
import { IconPhoto, IconUpload, IconX } from '@tabler/icons-react'
import { ReactNode, useCallback } from 'react'

const IMAGE_TYPES = ['image/png', 'image/jpeg']

function DropzoneText() {
  return (
    <div>
      <Text fw={700} size="sm" inline>
        Drag product image here
      </Text>
      <Text size="xs" color="dimmed" inline mt={7}>
        png, jpeg up to 1MB
      </Text>
    </div>
  )
}

interface ProductImageDropzoneProps {
  label: ReactNode
  name: string
  loading: boolean
  imageWidth: number
  imageHeight: number
  uploadedImageUrl?: string
  onDrop?: (file: FileWithPath) => void
  onRemove?: () => void
}

export function ProductImageDropzone({
  label,
  imageWidth,
  imageHeight,
  uploadedImageUrl,
  onDrop,
  onRemove,
  ...props
}: ProductImageDropzoneProps) {
  const theme = useMantineTheme()

  const handleDrop = useCallback(
    (files: FileWithPath[]) => {
      console.log('on drop!:', files)
      if (files.length > 0) {
        onDrop?.(files[0])
      }
    },
    [onDrop]
  )

  const DropzoneWrapper = useCallback(
    ({ children }: { children: ReactNode }) => {
      return (
        <>
          {uploadedImageUrl ? (
            <Indicator
              label={
                <UnstyledButton
                  onClick={onRemove}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyItems: 'center',
                  }}
                >
                  <IconX size="1rem" stroke={4} color="white" />
                </UnstyledButton>
              }
              size={24}
              styles={{ indicator: { padding: 0 } }}
              color="red"
            >
              {children}
            </Indicator>
          ) : (
            <div>{children}</div>
          )}
        </>
      )
    },
    [uploadedImageUrl, onRemove]
  )

  return (
    <Input.Wrapper label={label}>
      <div
        style={{
          width: `${imageWidth + 20}px`,
          height: `${imageHeight + 20}px`,
        }}
      >
        <DropzoneWrapper>
          <Dropzone
            styles={{ inner: { height: '100%' } }}
            accept={IMAGE_TYPES}
            multiple={false}
            maxSize={1 * 1024 ** 2}
            onDrop={handleDrop}
            w={imageWidth + 20}
            h={imageHeight + 20}
            {...props}
          >
            <Center h="100%">
              <Dropzone.Accept>
                <Group spacing="xs">
                  <IconUpload
                    size="3rem"
                    stroke={1.5}
                    color={theme.colors[theme.primaryColor][6]}
                  />
                  <DropzoneText />
                </Group>
              </Dropzone.Accept>
              <Dropzone.Reject>
                <Group spacing="xs">
                  <IconX size="3rem" stroke={1.5} color={theme.colors.red[6]} />
                  <DropzoneText />
                </Group>
              </Dropzone.Reject>
              <Dropzone.Idle>
                {uploadedImageUrl ? (
                  <FixedSizeImage
                    src={uploadedImageUrl}
                    alt="Product image"
                    width={imageWidth}
                    height={imageHeight}
                  />
                ) : (
                  <Group spacing="xs">
                    <IconPhoto size="3rem" stroke={1.5} />
                    <DropzoneText />
                  </Group>
                )}
              </Dropzone.Idle>
            </Center>
          </Dropzone>
        </DropzoneWrapper>
      </div>
    </Input.Wrapper>
  )
}
