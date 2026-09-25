import http from '@/utils/request'
import type { User } from '@/types'

export interface TokenResult {
  token: string
  user: User
}

export function sendCode(email: string) {
  return http.post<{ email: string; expire_seconds: number; dev_code?: string }>('/auth/send-code', { email })
}

export function register(data: { student_no: string; email: string; password: string; code: string }) {
  return http.post<TokenResult>('/auth/register', data)
}

export function login(data: { account: string; password: string }) {
  return http.post<TokenResult>('/auth/login', data)
}
