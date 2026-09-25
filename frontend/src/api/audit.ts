import http from '@/utils/request'
import type { AuditLog, PageData } from '@/types'

export function listAuditLogs(params: { page?: number; page_size?: number; action?: string; user_id?: number }) {
  return http.get<PageData<AuditLog>>('/audit-logs', params)
}
