import { Group, Text, TextProps, useMantineTheme } from '@mantine/core'
import { IconAlertCircle } from '@tabler/icons-react'
import { FC } from 'react'

type ErrorMessageProps = TextProps & { children?: string }

export const ErrorMessage: FC<ErrorMessageProps> = ({ children, ...rest }) => {
  const theme = useMantineTheme()
  if (!children?.length) return null
  return (
    <Text
      weight={500}
      size="sm"
      style={{
        wordBreak: 'break-word',
        display: 'block',
        position: 'relative',
      }}
      {...rest}
    >
      <Group spacing={5} sx={{ position: 'absolute' }}>
        <IconAlertCircle width={theme.fontSizes.lg} />
        {children}
      </Group>
    </Text>
  )
}
