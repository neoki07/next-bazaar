import { NextPage } from 'next'

export interface PageAuthConfig {
  auth?: boolean
}

export type Page<Props = {}, InitialProps = Props> = NextPage<
  Props,
  InitialProps
> &
  PageAuthConfig
