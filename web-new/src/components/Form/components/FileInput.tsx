import { Center, Group, FileInput as MantineFileInput } from '@mantine/core'
import { IconFileText, IconPhoto, IconUpload } from '@tabler/icons-react'
import { useController } from 'react-hook-form'
import { FileInputProps } from '../types'
import { ErrorMessage } from './ErrorMessage'

interface ValueProps {
  file: File | null
}

function Value({ file }: ValueProps) {
  if (!file) {
    return null
  }

  const Icon = file.type.includes('image') ? IconPhoto : IconFileText

  return (
    <Center
      inline
      sx={(theme) => ({
        backgroundColor:
          theme.colorScheme === 'dark'
            ? theme.colors.dark[7]
            : theme.colors.gray[1],
        fontSize: theme.fontSizes.xs,
        padding: '3px 7px',
        borderRadius: theme.radius.sm,
      })}
    >
      <Icon size={14} style={{ marginRight: 5 }} />
      <span
        style={{
          whiteSpace: 'nowrap',
          textOverflow: 'ellipsis',
          overflow: 'hidden',
          maxWidth: 200,
          display: 'inline-block',
        }}
      >
        {file.name}
      </span>
    </Center>
  )
}

interface ValueComponentProps {
  value: File | File[] | null
}

function ValueComponent({ value }: ValueComponentProps) {
  if (Array.isArray(value)) {
    return (
      <Group spacing="sm" py="xs">
        {value.map((file, index) => (
          <Value file={file} key={index} />
        ))}
      </Group>
    )
  }

  return <Value file={value} />
}

export function FileInput({ label, name, ...rest }: FileInputProps<boolean>) {
  const {
    field,
    fieldState: { error: fieldError },
  } = useController({ name })

  const error = fieldError ? (
    <ErrorMessage>{fieldError.message?.toString()}</ErrorMessage>
  ) : undefined

  return (
    <MantineFileInput
      label={label}
      icon={<IconUpload size={14} />}
      valueComponent={ValueComponent}
      error={error}
      {...rest}
      {...field}
    />
  )
}
