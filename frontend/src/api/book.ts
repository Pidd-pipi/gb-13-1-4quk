import http from '@/utils/request'
import type { Book, PageData } from '@/types'

export interface BookQuery {
  keyword?: string
  subject_category?: string
  condition?: string
  trade_type?: string
  min_price?: number
  max_price?: number
  seller_id?: number
  status?: string
  sort?: string
  page?: number
  page_size?: number
}

export interface BookPayload {
  title: string
  author: string
  isbn: string
  course_name: string
  original_price: number
  price: number
  condition: string
  subject_category: string
  trade_type: string
  campus: string
  description: string
  images: string[]
  lendable: boolean
  lend_days: number
}

export function listBooks(params: BookQuery) {
  return http.get<PageData<Book>>('/books', params)
}

export function getBook(id: number) {
  return http.get<Book>(`/books/${id}`)
}

export function createBook(data: BookPayload) {
  return http.post<Book>('/books', data)
}

export function updateBook(id: number, data: Partial<BookPayload>) {
  return http.put<Book>(`/books/${id}`, data)
}

export function deleteBook(id: number) {
  return http.delete<{ id: number }>(`/books/${id}`)
}

export function reserveBook(id: number) {
  return http.post<Book>(`/books/${id}/reserve`)
}

export function cancelReserveBook(id: number) {
  return http.post<Book>(`/books/${id}/cancel-reserve`)
}

export function markSold(id: number) {
  return http.post<Book>(`/books/${id}/sold`)
}

export function addFavorite(id: number) {
  return http.post<{ id: number; favorited: boolean }>(`/books/${id}/favorite`)
}

export function removeFavorite(id: number) {
  return http.delete<{ id: number; favorited: boolean }>(`/books/${id}/favorite`)
}

export function listFavorites(params: { page?: number; page_size?: number }) {
  return http.get<PageData<Book>>('/books/favorites', params)
}

export function listHistory(params: { limit?: number }) {
  return http.get<PageData<Book>>('/books/history', params)
}

export function getRecommendations(params: { limit?: number }) {
  return http.get<Book[]>('/books/recommendations', params)
}
