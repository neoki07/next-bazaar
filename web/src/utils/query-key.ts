import { QueryKey } from '@tanstack/react-query'

const QUERY_KEY_NON_CREDENTIALS = 'non-credentials'

/**
 * Add non-credentials to query key
 */
export function addNonCredentialsToQueryKey(queryKey: QueryKey): QueryKey {
  return [QUERY_KEY_NON_CREDENTIALS, ...queryKey]
}

/**
 * Returns true if query key is non-credentials
 */
export function isNonCredentialsQueryKey(queryKey: QueryKey) {
  return queryKey[0] === QUERY_KEY_NON_CREDENTIALS
}
