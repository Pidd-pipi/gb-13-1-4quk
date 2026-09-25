import http from '@/utils/request'
import type { Evaluation, PageData } from '@/types'

export function createEvaluation(data: { to_user_id: number; book_id: number; type: string; content: string }) {
  return http.post<Evaluation>('/evaluations', data)
}

export function listUserEvaluations(userId: number, params: { page?: number; page_size?: number }) {
  return http.get<PageData<Evaluation>>(`/users/${userId}/evaluations`, params)
}
