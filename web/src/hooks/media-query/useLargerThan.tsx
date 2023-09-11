import { MantineNumberSize, useMantineTheme } from '@mantine/core'
import { useMediaQuery } from '@mantine/hooks'

export function useLargerThan(breakpoint: MantineNumberSize) {
  const theme = useMantineTheme()
  return useMediaQuery(theme.fn.largerThan(breakpoint).slice(7))
}
