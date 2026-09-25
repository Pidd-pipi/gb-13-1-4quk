import http from '@/utils/request'
import type { Borrow, Conversation, PageData } from '@/types'

export interface BorrowQuery {
  role?: 'lender' | 'borrower'
  book_id?: number
  status?: string
  page?: number
  page_size?: number
}

// 某本书的借阅申请列表（卖家看全部 / 借阅人看自己）
export function listBookBorrows(bookId: number, params?: Omit<BorrowQuery, 'book_id'>) {
  return http.get<PageData<Borrow>>(`/books/${bookId}/borrows`, params)
}

// 我的借阅：role=lender 我借出的 / role=borrower 我借入的
export function listMyBorrows(params: BorrowQuery) {
  return http.get<PageData<Borrow>>('/borrows', params)
}

export function getBorrow(id: number) {
  return http.get<Borrow>(`/borrows/${id}`)
}

// 同学在书籍详情提交借阅申请
export function applyBorrow(bookId: number, message?: string) {
  return http.post<Borrow>(`/books/${bookId}/borrows`, { message: message || '' })
}

// 卖家同意借阅（书籍变为借出中，记录到期日，其他申请自动拒绝）
export function approveBorrow(id: number) {
  return http.post<Borrow>(`/borrows/${id}/approve`)
}

// 卖家拒绝借阅申请
export function rejectBorrow(id: number, reason?: string) {
  return http.post<Borrow>(`/borrows/${id}/reject`, { reason: reason || '' })
}

// 借阅人到期前归还（等待卖家确认）
export function returnBorrow(id: number) {
  return http.post<Borrow>(`/borrows/${id}/return`)
}

// 卖家确认拿回，书籍恢复可借
export function confirmReturnBorrow(id: number) {
  return http.post<Borrow>(`/borrows/${id}/confirm-return`)
}

// 卖家对逾期借阅发送提醒
export function remindBorrow(id: number, content?: string) {
  return http.post<Borrow>(`/borrows/${id}/remind`, { content: content || '' })
}

// 借阅双方沟通（卖家联系借阅人 / 借阅人联系卖家均可）
export function contactBorrow(id: number, content?: string) {
  return http.post<Conversation>(`/borrows/${id}/contact`, { content: content || '' })
}
