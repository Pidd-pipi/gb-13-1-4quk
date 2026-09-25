// 与后端 internal/constants/enums.go 对应的业务枚举。
// 新增枚举值必须同步修改：后端 constants/enums.go、formatters.go、错误码、日志模板、
// 前端本文件、筛选组件、状态徽标组件。

// 用户角色
export const Role = {
  STUDENT: 'student',
  ADMIN: 'admin',
} as const
export type RoleType = (typeof Role)[keyof typeof Role]

export const RoleText: Record<string, string> = {
  [Role.STUDENT]: '学生',
  [Role.ADMIN]: '管理员',
}

// 书籍状态机：on_sale(在售) -> reserved(已预约) -> sold(已售出)；
// 短借：on_sale <-> lent_out(借出，到期未还时显示逾期)
export const BookStatus = {
  ON_SALE: 'on_sale',
  RESERVED: 'reserved',
  SOLD: 'sold',
  LENT_OUT: 'lent_out',
} as const
export type BookStatusType = (typeof BookStatus)[keyof typeof BookStatus]

export const BookStatusText: Record<string, string> = {
  [BookStatus.ON_SALE]: '在售',
  [BookStatus.RESERVED]: '已预约',
  [BookStatus.SOLD]: '已售出',
  [BookStatus.LENT_OUT]: '借出',
}

export const BookStatusBadge: Record<string, 'success' | 'warning' | 'default' | 'primary'> = {
  [BookStatus.ON_SALE]: 'success',
  [BookStatus.RESERVED]: 'warning',
  [BookStatus.SOLD]: 'default',
  [BookStatus.LENT_OUT]: 'primary',
}

// 短借借期（发布时可借选项）
export const LendDays = {
  SHORT: 7,
  LONG: 14,
} as const

export const LendDaysOptions = [
  { value: LendDays.SHORT, label: '7 天' },
  { value: LendDays.LONG, label: '14 天' },
]

// 借阅申请状态机：pending(待同意) -> approved(借出中) ->
// returned(待确认归还) -> completed(已归还)；pending 也可 rejected(已拒绝)
export const BorrowStatus = {
  PENDING: 'pending',
  APPROVED: 'approved',
  REJECTED: 'rejected',
  RETURNED: 'returned',
  COMPLETED: 'completed',
} as const
export type BorrowStatusType = (typeof BorrowStatus)[keyof typeof BorrowStatus]

export const BorrowStatusText: Record<string, string> = {
  [BorrowStatus.PENDING]: '待同意',
  [BorrowStatus.APPROVED]: '借出中',
  [BorrowStatus.REJECTED]: '已拒绝',
  [BorrowStatus.RETURNED]: '待确认归还',
  [BorrowStatus.COMPLETED]: '已归还',
}

// 新旧程度
export const Condition = {
  BRAND_NEW: 'brand_new',
  NINE_NEW: 'nine_new',
  SEVEN_NEW: 'seven_new',
  FIVE_NEW: 'five_new',
} as const
export type ConditionType = (typeof Condition)[keyof typeof Condition]

export const ConditionText: Record<string, string> = {
  [Condition.BRAND_NEW]: '全新',
  [Condition.NINE_NEW]: '九成新',
  [Condition.SEVEN_NEW]: '七成新',
  [Condition.FIVE_NEW]: '五成新',
}

export const ConditionOptions = [
  { value: Condition.BRAND_NEW, label: ConditionText[Condition.BRAND_NEW] },
  { value: Condition.NINE_NEW, label: ConditionText[Condition.NINE_NEW] },
  { value: Condition.SEVEN_NEW, label: ConditionText[Condition.SEVEN_NEW] },
  { value: Condition.FIVE_NEW, label: ConditionText[Condition.FIVE_NEW] },
]

// 学科分类
export const SubjectCategory = {
  SCIENCE: 'science',
  HUMANITIES: 'humanities',
  ECON_MANAGEMENT: 'econ_management',
  ART: 'art',
  OTHER: 'other',
} as const
export type SubjectCategoryType = (typeof SubjectCategory)[keyof typeof SubjectCategory]

export const SubjectCategoryText: Record<string, string> = {
  [SubjectCategory.SCIENCE]: '理工',
  [SubjectCategory.HUMANITIES]: '文史',
  [SubjectCategory.ECON_MANAGEMENT]: '经管',
  [SubjectCategory.ART]: '艺术',
  [SubjectCategory.OTHER]: '其他',
}

export const SubjectCategoryOptions = [
  { value: SubjectCategory.SCIENCE, label: SubjectCategoryText[SubjectCategory.SCIENCE] },
  { value: SubjectCategory.HUMANITIES, label: SubjectCategoryText[SubjectCategory.HUMANITIES] },
  { value: SubjectCategory.ECON_MANAGEMENT, label: SubjectCategoryText[SubjectCategory.ECON_MANAGEMENT] },
  { value: SubjectCategory.ART, label: SubjectCategoryText[SubjectCategory.ART] },
  { value: SubjectCategory.OTHER, label: SubjectCategoryText[SubjectCategory.OTHER] },
]

// 交易方式
export const TradeType = {
  IN_PERSON: 'in_person',
  MAIL: 'mail',
} as const
export type TradeTypeType = (typeof TradeType)[keyof typeof TradeType]

export const TradeTypeText: Record<string, string> = {
  [TradeType.IN_PERSON]: '面交',
  [TradeType.MAIL]: '邮寄',
}

// 评价类型
export const EvaluationType = {
  GOOD: 'good',
  NEUTRAL: 'neutral',
  BAD: 'bad',
} as const
export type EvaluationTypeType = (typeof EvaluationType)[keyof typeof EvaluationType]

export const EvaluationTypeText: Record<string, string> = {
  [EvaluationType.GOOD]: '好评',
  [EvaluationType.NEUTRAL]: '中评',
  [EvaluationType.BAD]: '差评',
}

export const EvaluationTypeOptions = [
  { value: EvaluationType.GOOD, label: EvaluationTypeText[EvaluationType.GOOD] },
  { value: EvaluationType.NEUTRAL, label: EvaluationTypeText[EvaluationType.NEUTRAL] },
  { value: EvaluationType.BAD, label: EvaluationTypeText[EvaluationType.BAD] },
]

// 求购状态
export const WishStatus = {
  OPEN: 'open',
  CLOSED: 'closed',
} as const
export type WishStatusType = (typeof WishStatus)[keyof typeof WishStatus]

export const WishStatusText: Record<string, string> = {
  [WishStatus.OPEN]: '进行中',
  [WishStatus.CLOSED]: '已关闭',
}

// 排序方式
export const SortOptions = [
  { value: 'newest', label: '最新发布' },
  { value: 'price_asc', label: '价格最低' },
  { value: 'price_desc', label: '价格最高' },
  { value: 'most_viewed', label: '最多浏览' },
]
