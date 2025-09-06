import type { AxiosError } from 'axios'

export interface ApiError {
  errors?: { code: number; message: string }[]
  meta?: { http_status: number }
}

export function extractErrorMessage(error: unknown): string {
  const axError = error as AxiosError<ApiError>
  const apiError = axError.response?.data
  return apiError?.errors?.[0]?.message ?? 'Terjadi kesalahan!'
}
