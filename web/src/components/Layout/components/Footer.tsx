import { Anchor, Center, Footer as MantineFooter, Text } from '@mantine/core'

export function Footer() {
  return (
    <MantineFooter height={60} style={{ position: 'static' }}>
      <Center h="100%">
        <Text size="sm" sx={(theme) => ({ color: theme.colors.gray[7] })}>
          Built by ot07. The source code is available on{' '}
          <Anchor href="https://github.com/ot07/next-bazaar" target="_blank">
            GitHub
          </Anchor>
          .
        </Text>
      </Center>
    </MantineFooter>
  )
}
