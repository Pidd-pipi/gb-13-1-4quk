<template>
  <div class="page detail-page">
    <van-nav-bar title="书籍详情" left-arrow fixed placeholder @click-left="$router.back()" />

    <template v-if="book">
      <van-swipe v-if="images.length" class="detail-swipe" :autoplay="0" indicator-color="white">
        <van-swipe-item v-for="img in images" :key="img">
          <img :src="img" class="detail-img" alt="book" />
        </van-swipe-item>
      </van-swipe>
      <div v-else class="no-img"><van-icon name="book-o" size="48" color="#c8c9cc" /></div>

      <div class="detail-card">
        <div class="detail-title-row">
          <div class="detail-title">{{ book.title }}</div>
          <van-tag v-if="book.borrowable" plain type="primary" size="medium">可借 {{ book.borrow_duration }} 天</van-tag>
        </div>
        <div class="detail-price">
          <span class="price">{{ formatPrice(book.price) }}</span>
          <span v-if="book.original_price" class="origin">原价 {{ formatPrice(book.original_price) }}</span>
        </div>
        <div class="detail-tags">
          <ConditionTag :condition="book.condition" />
          <StatusBadge :status="book.status" />
          <van-tag plain type="primary">{{ book.subject_text }}</van-tag>
          <van-tag plain>{{ book.trade_type_text }}</van-tag>
        </div>

        <!-- 借阅进度条：到期日 / 逾期提示 -->
        <div v-if="activeBorrow" class="borrow-banner" :class="{ overdue: activeBorrow.is_overdue }">
          <template v-if="activeBorrow.is_overdue">
            <van-icon name="warning-o" />
            <span>已逾期，应于 {{ activeBorrow.due_at }} 归还</span>
          </template>
          <template v-else>
            <van-icon name="clock-o" />
            <span>{{ activeBorrow.status === 'returning' ? '借阅人已发起归还，等待卖家确认' : `借出中，到期日 ${activeBorrow.due_at}` }}</span>
          </template>
        </div>

        <van-cell-group inset class="meta-group">
          <van-cell title="作者" :value="book.author || '佚名'" />
          <van-cell title="ISBN" :value="book.isbn || '-'" />
          <van-cell title="课程" :value="book.course_name || '-'" />
          <van-cell title="校区" :value="book.campus || '-'" />
          <van-cell title="浏览量" :value="String(book.view_count)" />
          <van-cell v-if="book.borrowable" title="短借时长" :value="`${book.borrow_duration} 天`" />
        </van-cell-group>
        <div class="desc">
          <div class="desc-title">书籍描述</div>
          <div class="desc-body">{{ book.description || '暂无描述' }}</div>
        </div>
        <van-cell-group inset>
          <van-cell :title="book.seller?.name || '卖家'" :label="`${book.seller?.department || ''} · ${book.seller?.campus || ''}`" @click="goSeller">
            <template #icon><van-image round width="36" height="36" :src="book.seller?.avatar_url || fallbackAvatar" /></template>
            <template #value><van-icon name="arrow" /></template>
          </van-cell>
        </van-cell-group>

        <!-- 卖家视角：当前借阅人信息 -->
        <div v-if="isOwner && activeBorrow" class="borrower-card" @click="goProfile(activeBorrow.borrower_id)">
          <van-image round width="40" height="40" :src="activeBorrow.borrower?.avatar_url || fallbackAvatar" />
          <div class="borrower-info">
            <div class="borrower-name">{{ activeBorrow.borrower?.name || '借阅人' }}</div>
            <div class="borrower-meta">
              {{ activeBorrow.borrower?.department || '' }} · {{ activeBorrow.borrower?.contact || '未留联系方式' }}
            </div>
          </div>
          <van-tag v-if="activeBorrow.is_overdue" type="danger" size="medium">逾期</van-tag>
          <van-icon name="arrow" />
        </div>

        <!-- 卖家视角：待处理申请入口 -->
        <div v-if="isOwner && book.pending_borrow_count" class="pending-entry" @click="goBookBorrows">
          <van-icon name="envelop-o" />
          <span>有 {{ book.pending_borrow_count }} 条待处理的借阅申请</span>
          <van-icon name="arrow" />
        </div>

        <!-- 借阅人视角：我的申请状态 -->
        <van-notice-bar
          v-if="!isOwner && myBorrow"
          :type="myBorrow.status === 'rejected' ? 'danger' : myBorrow.is_overdue ? 'warning' : 'info'"
          :text="myBorrowStatusText"
        />
      </div>

      <div class="action-bar">
        <van-button plain type="warning" icon="star-o" @click="onFavorite">
          {{ book.is_favorite ? '已收藏' : '收藏' }}
        </van-button>
        <template v-if="isOwner">
          <van-button v-if="book.status === 'on_sale'" plain type="primary" icon="edit" @click="$router.push({ name: 'publish-book', query: { id: String(book.id) } })">编辑</van-button>
          <van-button v-if="book.status === 'on_sale'" type="danger" plain icon="delete-o" @click="onDelete">下架</van-button>
          <van-button v-if="book.status === 'on_sale'" type="warning" icon="records-o" @click="goBookBorrows">
            借阅申请<template v-if="book.pending_borrow_count">({{ book.pending_borrow_count }})</template>
          </van-button>
          <van-button v-if="book.status === 'reserved'" type="success" icon="passed" @click="onSold">确认售出</van-button>
          <van-button v-if="book.status === 'reserved'" plain type="default" @click="onCancelReserve">取消预约</van-button>
          <!-- 借出中：卖家确认收回 / 逾期提醒 / 联系借阅人 -->
          <template v-if="activeBorrow">
            <van-button v-if="activeBorrow.status === 'returning'" type="success" icon="passed" @click="onConfirmReturn">确认收回</van-button>
            <van-button v-if="activeBorrow.is_overdue" type="danger" plain icon="bell" @click="onRemind">发送提醒</van-button>
            <van-button plain type="primary" icon="chat-o" @click="onContactBorrower">联系借阅人</van-button>
          </template>
        </template>
        <template v-else>
          <van-button v-if="book.status === 'on_sale'" type="primary" icon="chat-o" @click="onContact">联系卖家</van-button>
          <van-button v-if="book.status === 'on_sale'" type="warning" plain icon="clock-o" @click="onReserve">预约</van-button>
          <!-- 可借且无本人进行中借阅：提交借阅申请 -->
          <van-button
            v-if="book.status === 'on_sale' && book.borrowable && !hasActiveMine && myBorrow?.status !== 'pending'"
            type="success"
            icon="exchange"
            @click="onApply"
          >
            申请借阅{{ book.borrow_duration ? ` ${book.borrow_duration} 天` : '' }}
          </van-button>
          <van-button v-if="myBorrow?.status === 'pending'" plain type="warning" icon="clock-o">申请待同意</van-button>
          <!-- 已被借出（非本人借阅）：提示不可借 -->
          <van-button v-if="book.status === 'loaned' && !hasActiveMine" plain type="default" disabled icon="clock-o">借出中</van-button>
          <!-- 我是当前借阅人：到期前归还 -->
          <van-button
            v-if="activeBorrow && activeBorrow.borrower_id === authId && activeBorrow.status === 'approved'"
            type="success"
            icon="revoke"
            @click="onReturn"
          >
            我要归还
          </van-button>
          <van-button v-if="book.status === 'reserved'" plain type="default" icon="clock-o" @click="onCancelReserve">取消预约</van-button>
        </template>
      </div>
    </template>
    <EmptyState v-else description="书籍不存在或已删除" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showConfirmDialog, showToast } from 'vant'
import { getBook, addFavorite, removeFavorite, reserveBook, cancelReserveBook, markSold, deleteBook } from '@/api/book'
import { applyBorrow, returnBorrow, confirmReturnBorrow, remindBorrow, contactBorrow } from '@/api/borrow'
import { createConversationFromBook } from '@/api/conversation'
import type { Book } from '@/types'
import { formatPrice, formatImages } from '@/utils/format'
import { useAuthStore } from '@/stores/authStore'
import ConditionTag from '@/components/ConditionTag.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import EmptyState from '@/components/EmptyState.vue'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const book = ref<Book | null>(null)
const fallbackAvatar = ''

const images = computed(() => (book.value ? formatImages(book.value.images) : []))
const isOwner = computed(() => !!book.value && !!auth.user && book.value.seller_id === auth.user.id)
const authId = computed(() => auth.user?.id || 0)
const activeBorrow = computed(() => book.value?.active_borrow || null)
const myBorrow = computed(() => book.value?.my_borrow || null)
const hasActiveMine = computed(
  () => !!activeBorrow.value && activeBorrow.value.borrower_id === authId.value,
)

const myBorrowStatusText = computed(() => {
  const b = myBorrow.value
  if (!b) return ''
  if (b.status === 'pending') return '借阅申请已提交，等待卖家同意'
  if (b.status === 'rejected') return `借阅申请被拒绝${b.reject_reason ? `：${b.reject_reason}` : ''}`
  if (b.status === 'returned') return `该笔借阅已归还（${b.confirmed_at}）`
  if (b.status === 'returning') return '你已发起归还，等待卖家确认收回'
  if (b.is_overdue) return `借阅已逾期，到期日 ${b.due_at}，请尽快归还`
  if (b.status === 'approved') return `借阅中，到期日 ${b.due_at}`
  return b.status_text
})

onMounted(async () => {
  await reload()
})

async function reload() {
  const id = Number(route.params.id)
  try {
    book.value = await getBook(id)
  } catch {
    book.value = null
  }
}

async function onFavorite() {
  if (!auth.isLoggedIn) {
    router.push('/login')
    return
  }
  if (!book.value) return
  try {
    if (book.value.is_favorite) {
      await removeFavorite(book.value.id)
      showToast('已取消收藏')
    } else {
      await addFavorite(book.value.id)
      showToast('收藏成功')
    }
    book.value.is_favorite = !book.value.is_favorite
    book.value.favorite_count = Math.max(0, (book.value.favorite_count || 0) + (book.value.is_favorite ? 1 : -1))
  } catch {
    /* toast handled by interceptor */
  }
}

async function onApply() {
  if (!auth.isLoggedIn) {
    router.push('/login')
    return
  }
  if (!book.value) return
  try {
    await showConfirmDialog({
      title: '申请借阅',
      message: `确认申请借阅《${book.value.title}》${book.value.borrow_duration} 天？卖家同意后开始计算到期日。`,
    })
    await applyBorrow(book.value.id)
    showToast('借阅申请已提交，等待卖家同意')
    await reload()
  } catch {
    /* cancelled or error */
  }
}

async function onReturn() {
  const b = activeBorrow.value
  if (!b) return
  try {
    await showConfirmDialog({ title: '归还书籍', message: '确认已把书交还给卖家？提交后等待卖家确认收回。' })
    await returnBorrow(b.id)
    showToast('已发起归还，等待卖家确认')
    await reload()
  } catch {
    /* cancelled */
  }
}

async function onConfirmReturn() {
  const b = activeBorrow.value
  if (!b) return
  try {
    await showConfirmDialog({ title: '确认收回', message: '确认已经拿回这本书？确认后书籍将恢复可借。' })
    await confirmReturnBorrow(b.id)
    showToast('已确认收回，书籍恢复可借')
    await reload()
  } catch {
    /* cancelled */
  }
}

async function onRemind() {
  const b = activeBorrow.value
  if (!b) return
  try {
    await showConfirmDialog({
      title: '发送还书提醒',
      message: `将通过站内消息提醒 ${b.borrower?.name || '借阅人'} 尽快归还，确认发送？`,
    })
    await remindBorrow(b.id)
    showToast('提醒已发送')
    await reload()
  } catch {
    /* cancelled */
  }
}

async function onContactBorrower() {
  const b = activeBorrow.value
  if (!b) return
  try {
    const conv = await contactBorrow(b.id, '你好，关于《' + (book.value?.title || '') + '》的借阅，想和你沟通一下。')
    router.push(`/chat/${conv.id}`)
  } catch {
    /* toast handled */
  }
}

async function onReserve() {
  if (!auth.isLoggedIn) {
    router.push('/login')
    return
  }
  if (!book.value) return
  try {
    await showConfirmDialog({ title: '预约', message: '确认预约这本书吗？' })
    book.value = await reserveBook(book.value.id)
    showToast('预约成功')
  } catch {
    /* cancelled or error */
  }
}

async function onCancelReserve() {
  if (!book.value) return
  try {
    await showConfirmDialog({ title: '取消预约', message: '确认取消预约吗？' })
    book.value = await cancelReserveBook(book.value.id)
    showToast('已取消预约')
  } catch {
    /* cancelled */
  }
}

async function onSold() {
  if (!book.value) return
  try {
    await showConfirmDialog({ title: '确认售出', message: '确认该书籍已完成交易并售出？' })
    book.value = await markSold(book.value.id)
    showToast('已标记售出')
  } catch {
    /* cancelled */
  }
}

async function onDelete() {
  if (!book.value) return
  try {
    await showConfirmDialog({ title: '下架', message: '确认下架该书籍？' })
    await deleteBook(book.value.id)
    showToast('已下架')
    router.back()
  } catch {
    /* cancelled */
  }
}

async function onContact() {
  if (!auth.isLoggedIn) {
    router.push('/login')
    return
  }
  if (!book.value) return
  try {
    const conv = await createConversationFromBook(book.value.id, '你好，我对这本书感兴趣，还在吗？')
    router.push(`/chat/${conv.id}`)
  } catch {
    /* toast handled */
  }
}

function goSeller() {
  goProfile(book.value?.seller_id)
}

function goProfile(userId?: number) {
  if (userId) {
    router.push(`/profile?userId=${userId}`)
  }
}

function goBookBorrows() {
  if (!book.value) return
  router.push(`/borrows?bookId=${book.value.id}`)
}
</script>

<style scoped>
.detail-swipe {
  height: 260px;
}
.detail-img {
  width: 100%;
  height: 260px;
  object-fit: cover;
}
.no-img {
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f2f3f5;
}
.detail-card {
  background: #fff;
  border-radius: 8px;
  padding: 14px;
  margin: 12px 0;
}
.detail-title-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 8px;
}
.detail-title {
  font-size: 18px;
  font-weight: 700;
}
.detail-price {
  margin: 8px 0;
}
.price {
  color: #ee0a24;
  font-size: 22px;
  font-weight: 700;
}
.origin {
  margin-left: 8px;
  color: #969799;
  font-size: 13px;
  text-decoration: line-through;
}
.detail-tags {
  display: flex;
  gap: 6px;
  align-items: center;
  flex-wrap: wrap;
}
.borrow-banner {
  margin-top: 10px;
  padding: 8px 12px;
  border-radius: 6px;
  background: #ecf5ff;
  color: #1989fa;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 6px;
}
.borrow-banner.overdue {
  background: #ffeceb;
  color: #ee0a24;
}
.meta-group {
  margin: 12px 0;
}
.desc {
  margin-top: 12px;
}
.desc-title {
  font-weight: 600;
  margin-bottom: 6px;
}
.desc-body {
  color: #646566;
  font-size: 14px;
  line-height: 1.6;
  white-space: pre-wrap;
}
.borrower-card {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 12px;
  padding: 10px 12px;
  background: #f7f8fa;
  border-radius: 8px;
}
.borrower-info {
  flex: 1;
  min-width: 0;
}
.borrower-name {
  font-weight: 600;
  font-size: 14px;
}
.borrower-meta {
  color: #969799;
  font-size: 12px;
}
.pending-entry {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 12px;
  padding: 10px 12px;
  background: #fff7e8;
  color: #ff976a;
  border-radius: 8px;
  font-size: 13px;
}
.pending-entry span {
  flex: 1;
}
.action-bar {
  position: sticky;
  bottom: 0;
  display: flex;
  gap: 10px;
  padding: 10px 12px;
  background: #fff;
  border-top: 1px solid #ebedf0;
}
.action-bar .van-button {
  flex: 1;
}
</style>
