export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface PageData<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

export interface User {
  id: number
  student_no: string
  email: string
  name: string
  department: string
  campus: string
  contact: string
  avatar_url: string
  role: 'student' | 'admin'
  email_verified: boolean
  status: string
  created_at: string
}

export interface Book {
  id: number
  seller_id: number
  title: string
  author: string
  isbn: string
  course_name: string
  original_price: number
  price: number
  condition: string
  condition_text: string
  subject_category: string
  subject_text: string
  trade_type: string
  trade_type_text: string
  campus: string
  description: string
  images: string[]
  status: string
  status_text: string
  reserved_by: number
  view_count: number
  favorite_count: number
  created_at: string
  seller?: User
  is_favorite?: boolean
}

export interface Wish {
  id: number
  user_id: number
  book_title: string
  author: string
  isbn: string
  expected_price: number
  condition_requirement: string
  condition_text: string
  subject_category: string
  subject_text: string
  description: string
  status: string
  status_text: string
  created_at: string
  user?: User
}

export interface Conversation {
  id: number
  book_id: number
  wish_id: number
  buyer_id: number
  seller_id: number
  last_message: string
  last_message_at: string
  created_at: string
  book?: Book
  buyer?: User
  seller?: User
  unread_count: number
}

export interface Message {
  id: number
  conversation_id: number
  sender_id: number
  content: string
  image_url: string
  is_read: boolean
  created_at: string
  sender?: User
}

export interface Evaluation {
  id: number
  from_user_id: number
  to_user_id: number
  book_id: number
  type: string
  type_text: string
  content: string
  created_at: string
  from_user?: User
  book?: Book
}

export interface UserStats {
  user_id: number
  total_books: number
  on_sale_books: number
  sold_books: number
  total_evaluations: number
  good_count: number
  neutral_count: number
  bad_count: number
  good_rate: string
  risk_flagged: boolean
}

export interface AuditLog {
  id: number
  user_id: number
  action: string
  resource_type: string
  resource_id: number
  detail: string
  ip: string
  request_id: string
  created_at: string
}

export interface UploadResult {
  url: string
  object: string
  filename: string
  size: number
}
