// 共享格式化工具（对应后端 internal/util/formatters.go）
export function formatPrice(price: number): string {
  return `¥${Number(price || 0).toFixed(2)}`
}

export function formatTime(time?: string): string {
  if (!time) return ''
  return time
}

export function formatDate(time?: string): string {
  if (!time) return ''
  return time.slice(0, 10)
}

export function shortText(text?: string, max = 40): string {
  if (!text) return ''
  return text.length > max ? text.slice(0, max) + '…' : text
}

export function formatImages(images: string[] | string | undefined): string[] {
  if (!images) return []
  if (Array.isArray(images)) return images
  if (typeof images === 'string') {
    if (images === '[]' || images === '') return []
    try {
      const parsed = JSON.parse(images)
      return Array.isArray(parsed) ? parsed : []
    } catch {
      return images.split(',').filter(Boolean)
    }
  }
  return []
}
