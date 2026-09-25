import http from '@/utils/request'
import type { Conversation, Message, PageData } from '@/types'

export function listConversations(params: { page?: number; page_size?: number }) {
  return http.get<PageData<Conversation>>('/conversations', params)
}

export function getConversation(id: number) {
  return http.get<Conversation>(`/conversations/${id}`)
}

export function createConversationFromBook(bookId: number, content: string) {
  return http.post<Conversation>('/conversations', { book_id: bookId, content })
}

// 卖家主动联系某位同学（如借阅人）
export function createConversationToUser(bookId: number, toUserId: number, content: string) {
  return http.post<Conversation>('/conversations', { book_id: bookId, to_user_id: toUserId, content })
}

export function listMessages(conversationId: number) {
  return http.get<Message[]>(`/conversations/${conversationId}/messages`)
}

export function sendMessage(conversationId: number, content: string, imageUrl?: string) {
  return http.post<Message>(`/conversations/${conversationId}/messages`, { content, image_url: imageUrl })
}

export function markRead(conversationId: number) {
  return http.put<{ id: number }>(`/conversations/${conversationId}/read`)
}
