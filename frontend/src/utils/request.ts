import axios, { AxiosError, AxiosRequestConfig } from 'axios'
import { showToast } from 'vant'
import type { ApiResponse } from '@/types'

const TOKEN_KEY = 'campusbooks_token'

export function getToken(): string {
  return localStorage.getItem(TOKEN_KEY) || ''
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

const instance = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
})

instance.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 统一响应拦截：code !== 0 时提示错误并 reject
instance.interceptors.response.use(
  (response) => {
    const body = response.data as ApiResponse
    if (body.code !== 0) {
      showToast(body.message || '请求失败')
      return Promise.reject(new Error(body.message || '请求失败'))
    }
    return response
  },
  (error: AxiosError<ApiResponse>) => {
    const message = error.response?.data?.message || error.message || '网络异常'
    showToast(message)
    if (error.response?.status === 401) {
      clearToken()
      if (!window.location.pathname.startsWith('/login')) {
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  },
)

export async function request<T = unknown>(config: AxiosRequestConfig): Promise<T> {
  const res = await instance.request<ApiResponse<T>>(config)
  return res.data.data
}

export const http = {
  get: <T>(url: string, params?: object) =>
    request<T>({ method: 'GET', url, params }),
  post: <T>(url: string, data?: unknown) => request<T>({ method: 'POST', url, data }),
  put: <T>(url: string, data?: unknown) => request<T>({ method: 'PUT', url, data }),
  delete: <T>(url: string) => request<T>({ method: 'DELETE', url }),
}

export default http
