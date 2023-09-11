import { MantineNumberSize, useMantineTheme } from '@mantine/core'
import { useMediaQuery } from '@mantine/hooks'

export function useSmallerThan(breakpoint: MantineNumberSize) {
  const theme = useMantineTheme()
  return useMediaQuery(theme.fn.smallerThan(breakpoint).slice(7))
}
