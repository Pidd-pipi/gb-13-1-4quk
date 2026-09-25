import http from '@/utils/request'
import type { Evaluation, PageData, User, UserStats } from '@/types'

export function getMe() {
  return http.get<User>('/users/me')
}

export function updateProfile(data: { name: string; department: string; campus: string; contact: string }) {
  return http.put<User>('/users/me', data)
}

export function updateAvatar(avatarUrl: string) {
  return http.put<User>('/users/me/avatar', { avatar_url: avatarUrl })
}

export function getStats(userId?: number) {
  const url = userId ? `/users/${userId}/stats` : '/users/me/stats'
  return http.get<UserStats>(url)
}

export function listUsers(params: { page?: number; page_size?: number; role?: string }) {
  return http.get<PageData<User>>('/users', params)
}

export function getUserEvaluations(userId: number, params: { page?: number; page_size?: number }) {
  return http.get<PageData<Evaluation>>(`/users/${userId}/evaluations`, params)
}
