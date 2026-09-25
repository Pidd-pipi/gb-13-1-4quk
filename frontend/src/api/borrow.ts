import http from '@/utils/request'
import type { BorrowRequest, PageData } from '@/types'

// 同学在书籍详情提交借阅申请
export function applyBorrow(bookId: number) {
  return http.post<BorrowRequest>(`/books/${bookId}/borrow-requests`)
}

// 卖家查看某本书的全部借阅申请
export function listBookBorrows(bookId: number) {
  return http.get<BorrowRequest[]>(`/books/${bookId}/borrow-requests`)
}

// 我借入 / 我收到的借阅申请
export function listBorrows(role: 'borrower' | 'seller', params: { page?: number; page_size?: number } = {}) {
  return http.get<PageData<BorrowRequest>>('/borrow-requests', { role, ...params })
}

// 卖家同意
export function approveBorrow(id: number) {
  return http.post<BorrowRequest>(`/borrow-requests/${id}/approve`)
}

// 卖家拒绝
export function rejectBorrow(id: number) {
  return http.post<BorrowRequest>(`/borrow-requests/${id}/reject`)
}

// 借阅人归还
export function returnBorrow(id: number) {
  return http.post<BorrowRequest>(`/borrow-requests/${id}/return`)
}

// 卖家确认拿回
export function confirmReturnBorrow(id: number) {
  return http.post<BorrowRequest>(`/borrow-requests/${id}/confirm-return`)
}

// 卖家发送还书提醒
export function remindBorrow(id: number) {
  return http.post<BorrowRequest>(`/borrow-requests/${id}/remind`)
}
