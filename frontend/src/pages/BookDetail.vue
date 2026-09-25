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
        <div class="detail-title">{{ book.title }}</div>
        <div class="detail-price">
          <span class="price">{{ formatPrice(book.price) }}</span>
          <span v-if="book.original_price" class="origin">原价 {{ formatPrice(book.original_price) }}</span>
        </div>
        <div class="detail-tags">
          <ConditionTag :condition="book.condition" />
          <StatusBadge :status="book.status" />
          <van-tag plain type="primary">{{ book.subject_text }}</van-tag>
          <van-tag plain>{{ book.trade_type_text }}</van-tag>
          <van-tag v-if="book.lendable" plain type="success">可借 · {{ book.lend_days }} 天</van-tag>
        </div>
        <van-cell-group inset class="meta-group">
          <van-cell title="作者" :value="book.author || '佚名'" />
          <van-cell title="ISBN" :value="book.isbn || '-'" />
          <van-cell title="课程" :value="book.course_name || '-'" />
          <van-cell title="校区" :value="book.campus || '-'" />
          <van-cell title="浏览量" :value="String(book.view_count)" />
        </van-cell-group>
        <div class="desc">
          <div class="desc-title">书籍描述</div>
          <div class="desc-body">{{ book.description || '暂无描述' }}</div>
        </div>

        <template v-if="activeBorrow">
          <div class="borrow-title">
            借阅信息
            <van-tag v-if="activeBorrow.overdue" type="danger" round>已逾期</van-tag>
          </div>
          <van-cell-group inset class="borrow-group">
            <van-cell title="借阅状态" :value="activeBorrow.status_text" />
            <van-cell title="到期日" :value="activeBorrow.due_at || '-'">
              <template #value>
                <span :class="{ overdue: activeBorrow.overdue }">{{ activeBorrow.due_at || '-' }}</span>
              </template>
            </van-cell>
            <van-cell v-if="isOwner && activeBorrow.borrower" title="借阅人" is-link
              :label="`${activeBorrow.borrower.department || ''} · ${activeBorrow.borrower.contact || '未留联系方式'}`"
              :value="activeBorrow.borrower.name || '同学'"
              @click="goBorrower(activeBorrow.borrower_id)" />
          </van-cell-group>
        </template>

        <template v-if="isOwner && pendingBorrows.length">
          <div class="borrow-title">借阅申请（{{ pendingBorrows.length }}）</div>
          <van-cell-group inset class="borrow-group">
            <van-cell v-for="r in pendingBorrows" :key="r.id" :label="`申请于 ${r.created_at}`">
              <template #title>
                <span>{{ r.borrower?.name || '同学' }} 申请借阅 {{ book.lend_days }} 天</span>
              </template>
              <template #value>
                <div class="borrow-actions">
                  <van-button size="small" type="success" :loading="actingId === r.id" @click="onApproveBorrow(r)">同意</van-button>
                  <van-button size="small" plain type="default" :loading="actingId === r.id" @click="onRejectBorrow(r)">拒绝</van-button>
                </div>
              </template>
            </van-cell>
          </van-cell-group>
        </template>

        <van-cell-group inset>
          <van-cell :title="book.seller?.name || '卖家'" :label="`${book.seller?.department || ''} · ${book.seller?.campus || ''}`" @click="goSeller">
            <template #icon><van-image round width="36" height="36" :src="book.seller?.avatar_url || fallbackAvatar" /></template>
            <template #value><van-icon name="arrow" /></template>
          </van-cell>
        </van-cell-group>
      </div>

      <div class="action-bar">
        <van-button plain type="warning" icon="star-o" @click="onFavorite">
          {{ book.is_favorite ? '已收藏' : '收藏' }}
        </van-button>
        <template v-if="isOwner">
          <van-button v-if="book.status === 'on_sale'" plain type="primary" icon="edit" @click="$router.push({ name: 'publish-book', query: { id: String(book.id) } })">编辑</van-button>
          <van-button v-if="book.status === 'on_sale'" type="danger" icon="delete-o" @click="onDelete">下架</van-button>
          <van-button v-if="book.status === 'reserved'" type="success" icon="passed" @click="onSold">确认售出</van-button>
          <van-button v-if="book.status === 'reserved'" plain type="default" @click="onCancelReserve">取消预约</van-button>
          <template v-if="activeBorrow">
            <van-button plain type="warning" icon="bullhorn-o" :loading="acting" @click="onRemind">发送提醒</van-button>
            <van-button v-if="activeBorrow.status === 'returned'" type="success" icon="passed" :loading="acting" @click="onConfirmReturn">确认拿回</van-button>
          </template>
        </template>
        <template v-else>
          <van-button v-if="book.status === 'on_sale'" plain type="primary" icon="chat-o" @click="onContact">联系卖家</van-button>
          <van-button v-if="book.status === 'on_sale'" plain type="warning" icon="clock-o" @click="onReserve">预约</van-button>
          <van-button v-if="book.status === 'on_sale' && book.lendable && !myBorrow" type="primary" icon="records-o" @click="onApplyBorrow">
            申请借阅{{ book.lend_days ? ` ${book.lend_days} 天` : '' }}
          </van-button>
          <van-button v-if="myBorrow?.status === 'pending'" plain type="default" icon="clock-o" disabled>借阅待同意</van-button>
          <van-button v-if="myBorrow?.status === 'approved'" type="primary" icon="revoke" @click="onReturnBorrow">申请归还</van-button>
          <van-button v-if="myBorrow?.status === 'returned'" plain type="default" icon="passed" disabled>待卖家确认</van-button>
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
import {
  applyBorrow,
  approveBorrow,
  rejectBorrow,
  returnBorrow,
  confirmReturnBorrow,
  remindBorrow,
  listBookBorrows,
} from '@/api/borrow'
import { createConversationFromBook } from '@/api/conversation'
import type { Book, BorrowRequest } from '@/types'
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
const activeBorrow = computed(() => book.value?.active_borrow)
const myBorrow = computed(() => book.value?.my_borrow)
const pendingBorrows = ref<BorrowRequest[]>([])
const acting = ref(false)
const actingId = ref(0)

onMounted(reload)

async function reload() {
  const id = Number(route.params.id)
  try {
    book.value = await getBook(id)
    if (book.value && auth.user && book.value.seller_id === auth.user.id) {
      try {
        const list = await listBookBorrows(id)
        pendingBorrows.value = list.filter((r) => r.status === 'pending')
      } catch {
        pendingBorrows.value = []
      }
    }
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

async function onApplyBorrow() {
  if (!auth.isLoggedIn) {
    router.push('/login')
    return
  }
  if (!book.value) return
  try {
    await showConfirmDialog({
      title: '申请借阅',
      message: `申请借阅这本书 ${book.value.lend_days} 天？卖家同意后开始计时，请在到期前归还。`,
    })
    await applyBorrow(book.value.id)
    showToast('已提交借阅申请，等待卖家同意')
    await reload()
  } catch {
    /* cancelled */
  }
}

async function onApproveBorrow(r: BorrowRequest) {
  try {
    await showConfirmDialog({
      title: '同意借阅',
      message: `同意 ${r.borrower?.name || '该同学'} 借阅 ${book.value?.lend_days || 7} 天？同意后其他申请将自动拒绝。`,
    })
    actingId.value = r.id
    await approveBorrow(r.id)
    showToast('已同意，书籍标记为借出')
    await reload()
  } catch {
    /* cancelled */
  } finally {
    actingId.value = 0
  }
}

async function onRejectBorrow(r: BorrowRequest) {
  try {
    await showConfirmDialog({ title: '拒绝借阅', message: `拒绝 ${r.borrower?.name || '该同学'} 的借阅申请？` })
    actingId.value = r.id
    await rejectBorrow(r.id)
    showToast('已拒绝')
    await reload()
  } catch {
    /* cancelled */
  } finally {
    actingId.value = 0
  }
}

async function onReturnBorrow() {
  if (!myBorrow.value) return
  try {
    await showConfirmDialog({ title: '申请归还', message: '确认已将书归还给卖家？卖家确认拿回后借阅完成。' })
    acting.value = true
    await returnBorrow(myBorrow.value.id)
    showToast('已通知卖家，请等待确认拿回')
    await reload()
  } catch {
    /* cancelled */
  } finally {
    acting.value = false
  }
}

async function onConfirmReturn() {
  if (!activeBorrow.value) return
  try {
    await showConfirmDialog({ title: '确认拿回', message: '确认已拿回这本书？确认后书籍恢复可借。' })
    acting.value = true
    await confirmReturnBorrow(activeBorrow.value.id)
    showToast('已确认拿回，书籍恢复可借')
    await reload()
  } catch {
    /* cancelled */
  } finally {
    acting.value = false
  }
}

async function onRemind() {
  if (!activeBorrow.value) return
  try {
    const message = activeBorrow.value.overdue
      ? '书籍已逾期，确认向借阅人发送还书提醒？'
      : '确认向借阅人发送还书提醒？'
    await showConfirmDialog({ title: '发送提醒', message })
    acting.value = true
    await remindBorrow(activeBorrow.value.id)
    showToast('提醒已通过站内消息发送')
  } catch {
    /* cancelled */
  } finally {
    acting.value = false
  }
}

function goSeller() {
  if (book.value?.seller_id) {
    router.push(`/profile?userId=${book.value.seller_id}`)
  }
}

function goBorrower(userId: number) {
  if (userId) {
    router.push(`/profile?userId=${userId}`)
  }
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
.borrow-title {
  font-weight: 600;
  font-size: 14px;
  margin: 14px 0 6px;
  padding: 0 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.borrow-group {
  margin-bottom: 12px;
}
.borrow-group .overdue {
  color: #ee0a24;
  font-weight: 600;
}
.borrow-actions {
  display: flex;
  gap: 6px;
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
