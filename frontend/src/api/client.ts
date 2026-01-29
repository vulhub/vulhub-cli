import ky from 'ky'
import type { ApiResponse } from '@/types'

const api = ky.create({
  prefixUrl: '/api/v1',
  timeout: 60000,
  hooks: {
    afterResponse: [
      async (_request, _options, response) => {
        if (!response.ok) {
          const body = await response.json() as ApiResponse<unknown>
          if (body.error) {
            throw new Error(body.error.message)
          }
        }
        return response
      }
    ]
  }
})

export async function fetchApi<T>(
  endpoint: string,
  options?: Parameters<typeof api>[1]
): Promise<T> {
  const response = await api(endpoint, options).json<ApiResponse<T>>()
  if (!response.success || !response.data) {
    throw new Error(response.error?.message || 'Request failed')
  }
  return response.data
}

export async function postApi<T>(
  endpoint: string,
  body?: unknown
): Promise<T> {
  const response = await api.post(endpoint, { json: body }).json<ApiResponse<T>>()
  if (!response.success) {
    throw new Error(response.error?.message || 'Request failed')
  }
  return response.data as T
}

export async function deleteApi<T>(
  endpoint: string,
  body?: unknown
): Promise<T> {
  const response = await api.delete(endpoint, { json: body }).json<ApiResponse<T>>()
  if (!response.success) {
    throw new Error(response.error?.message || 'Request failed')
  }
  return response.data as T
}

export { api }
