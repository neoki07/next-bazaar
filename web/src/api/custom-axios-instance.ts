import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios'

export const AXIOS_INSTANCE = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
})

export const customAxiosInstance = <T>(
  config: AxiosRequestConfig,
  options?: AxiosRequestConfig
): Promise<AxiosResponse<T, any>> => {
  return AXIOS_INSTANCE.request({ ...config, ...options })
}

export type ErrorType<Error> = AxiosError<Error>
