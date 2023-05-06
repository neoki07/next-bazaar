import { FileInput as MantineFileInput } from '@mantine/core'
import { FileInputProps } from '../types'
import { IconUpload } from '@tabler/icons-react'
import { useController } from 'react-hook-form'
import { ErrorMessage } from './ErrorMessage'
import {
  FileInputProps as MantineFileInputProps,
  Group,
  Center,
} from '@mantine/core'
import {
  IconFileText,
  IconPhoto,
  Icon as TablerIcon,
} from '@tabler/icons-react'
import { FC } from 'react'

type ValueProps = {
  file: File | null
}

const Value: FC<ValueProps> = ({ file }) => {
  if (!file) return null
  let Icon: TablerIcon

  if (file.type.includes('image')) {
    Icon = IconPhoto
  } else {
    Icon = IconFileText
  }

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

const ValueComponent: MantineFileInputProps['valueComponent'] = ({ value }) => {
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

export const FileInput: FC<FileInputProps<boolean>> = ({
  label,
  name,
  ...rest
}) => {
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
