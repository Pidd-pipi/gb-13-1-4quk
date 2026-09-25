<template>
  <div class="page">
    <van-nav-bar title="短借管理" left-arrow fixed placeholder @click-left="$router.back()" />
    <van-tabs v-model:active="role" @change="reload">
      <van-tab title="我借出的" name="lender" />
      <van-tab title="我借入的" name="borrower" />
    </van-tabs>

    <van-pull-refresh v-model="refreshing" @refresh="reload">
      <div v-if="list.length" class="borrow-list">
        <div v-for="b in list" :key="b.id" class="borrow-card" @click="goDetail(b.book_id)">
          <div class="card-head">
            <span class="book-title van-multi-ellipsis--l2">{{ b.book?.title || `书籍 #${b.book_id}` }}</span>
            <van-tag :type="badgeType(b)">{{ statusLabel(b) }}</van-tag>
          </div>
          <div class="card-meta">
            <span v-if="role === 'lender'">借阅人：{{ b.borrower?.name || '同学' }}</span>
            <span v-else>出借人：{{ b.lender?.name || '卖家' }}</span>
            <span>借期 {{ b.duration }} 天</span>
          </div>
          <div class="card-time">
            <span>申请于 {{ b.created_at }}</span>
            <span v-if="b.due_at">到期日 {{ b.due_at }}</span>
          </div>
          <div v-if="b.status === 'rejected' && b.reject_reason" class="reject-reason">拒绝原因：{{ b.reject_reason }}</div>
          <div v-if="b.is_overdue" class="overdue-tip">
            <van-icon name="warning-o" /> 已超过到期日，请尽快归还
          </div>

          <div class="card-actions" @click.stop>
            <template v-if="role === 'lender'">
              <van-button v-if="b.status === 'pending'" type="success" size="small" @click="onApprove(b)">同意</van-button>
              <van-button v-if="b.status === 'pending'" plain type="danger" size="small" @click="onReject(b)">拒绝</van-button>
              <van-button v-if="b.status === 'returning'" type="primary" size="small" @click="onConfirm(b)">确认收回</van-button>
              <van-button v-if="b.is_overdue && b.status === 'approved'" type="warning" size="small" icon="bell" @click="onRemind(b)">提醒</van-button>
              <van-button v-if="b.is_overdue || b.status === 'approved' || b.status === 'returning'" plain size="small" @click="contact(b)">联系借阅人</van-button>
            </template>
            <template v-else>
              <van-button v-if="b.status === 'approved'" type="success" size="small" @click="onReturn(b)">我要归还</van-button>
              <van-button plain type="primary" size="small" @click="contact(b)">联系卖家</van-button>
            </template>
          </div>
        </div>
      </div>
      <EmptyState v-else description="暂无借阅记录" />
    </van-pull-refresh>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showConfirmDialog, showToast } from 'vant'
import { listMyBorrows, approveBorrow, rejectBorrow, returnBorrow, confirmReturnBorrow, remindBorrow, contactBorrow } from '@/api/borrow'
import { BorrowStatusBadge, BorrowStatusText } from '@/constants/enums'
import type { Borrow } from '@/types'
import EmptyState from '@/components/EmptyState.vue'

const route = useRoute()
const router = useRouter()
const role = ref<'lender' | 'borrower'>(String(route.query.role || 'lender') === 'borrower' ? 'borrower' : 'lender')
const bookId = ref<number | undefined>(route.query.bookId ? Number(route.query.bookId) : undefined)
const list = ref<Borrow[]>([])
const refreshing = ref(false)

onMounted(reload)

async function reload() {
  refreshing.value = true
  try {
    const data = await listMyBorrows({ role: role.value, book_id: bookId.value, page: 1, page_size: 50 })
    list.value = data.list
  } catch {
    list.value = []
  } finally {
    refreshing.value = false
  }
}

function statusLabel(b: Borrow): string {
  if (b.is_overdue) return '已逾期'
  return BorrowStatusText[b.status] || b.status
}

function badgeType(b: Borrow): 'success' | 'warning' | 'default' | 'primary' | 'danger' {
  if (b.is_overdue) return 'danger'
  return (BorrowStatusBadge[b.status] as 'success' | 'warning' | 'default' | 'primary' | 'danger') || 'default'
}

async function onApprove(b: Borrow) {
  try {
    await showConfirmDialog({
      title: '同意借阅',
      message: `同意将《${b.book?.title || ''}》借给 ${b.borrower?.name || '该同学'} ${b.duration} 天？同意后其他申请会被拒绝。`,
    })
    await approveBorrow(b.id)
    showToast('已同意，书籍标记为借出中')
    await reload()
  } catch {
    /* cancelled */
  }
}

async function onReject(b: Borrow) {
  try {
    await showConfirmDialog({ title: '拒绝借阅', message: '确认拒绝该借阅申请吗？' })
    await rejectBorrow(b.id)
    showToast('已拒绝')
    await reload()
  } catch {
    /* cancelled */
  }
}

async function onReturn(b: Borrow) {
  try {
    await showConfirmDialog({ title: '归还书籍', message: '确认已把书交还给卖家？提交后等待卖家确认收回。' })
    await returnBorrow(b.id)
    showToast('已发起归还，等待卖家确认')
    await reload()
  } catch {
    /* cancelled */
  }
}

async function onConfirm(b: Borrow) {
  try {
    await showConfirmDialog({ title: '确认收回', message: '确认已经拿回这本书？确认后书籍恢复可借。' })
    await confirmReturnBorrow(b.id)
    showToast('已确认收回')
    await reload()
  } catch {
    /* cancelled */
  }
}

async function onRemind(b: Borrow) {
  try {
    await showConfirmDialog({ title: '发送提醒', message: `将通过站内消息提醒 ${b.borrower?.name || '借阅人'} 尽快归还。` })
    await remindBorrow(b.id)
    showToast('提醒已发送')
    await reload()
  } catch {
    /* cancelled */
  }
}

async function contact(b: Borrow) {
  try {
    const conv = await contactBorrow(b.id, '你好，关于短借的事情想和你沟通一下。')
    router.push(`/chat/${conv.id}`)
  } catch {
    /* toast handled */
  }
}

function goDetail(id: number) {
  router.push(`/books/${id}`)
}
</script>

<style scoped>
.borrow-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 12px;
}
.borrow-card {
  background: #fff;
  border-radius: 8px;
  padding: 12px;
}
.card-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 8px;
}
.book-title {
  font-size: 15px;
  font-weight: 600;
  flex: 1;
  min-width: 0;
}
.card-meta,
.card-time {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 6px;
  color: #646566;
  font-size: 13px;
}
.reject-reason,
.overdue-tip {
  margin-top: 6px;
  font-size: 12px;
}
.reject-reason {
  color: #969799;
}
.overdue-tip {
  color: #ee0a24;
  display: flex;
  align-items: center;
  gap: 4px;
}
.card-actions {
  display: flex;
  gap: 8px;
  margin-top: 10px;
  justify-content: flex-end;
}
</style>
