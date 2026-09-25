import http from '@/utils/request'
import type { UploadResult } from '@/types'

export function uploadFile(file: File): Promise<UploadResult> {
  const form = new FormData()
  form.append('file', file)
  return http.post<UploadResult>('/uploads', form)
}
