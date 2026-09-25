import http from '@/utils/request'
import type { PageData, Wish } from '@/types'

export interface WishQuery {
  keyword?: string
  subject_category?: string
  user_id?: number
  status?: string
  page?: number
  page_size?: number
}

export interface WishPayload {
  book_title: string
  author: string
  isbn: string
  expected_price: number
  condition_requirement: string
  subject_category: string
  description: string
}

export function listWishes(params: WishQuery) {
  return http.get<PageData<Wish>>('/wishes', params)
}

export function getWish(id: number) {
  return http.get<Wish>(`/wishes/${id}`)
}

export function createWish(data: WishPayload) {
  return http.post<Wish>('/wishes', data)
}

export function updateWish(id: number, data: Partial<WishPayload>) {
  return http.put<Wish>(`/wishes/${id}`, data)
}

export function closeWish(id: number) {
  return http.post<Wish>(`/wishes/${id}/close`)
}

export function deleteWish(id: number) {
  return http.delete<{ id: number }>(`/wishes/${id}`)
}

export function contactWish(id: number, content: string) {
  return http.post<ConversationResult>(`/wishes/${id}/contact`, { wish_id: id, content })
}

export interface ConversationResult {
  id: number
  book_id: number
  wish_id: number
  buyer_id: number
  seller_id: number
}
