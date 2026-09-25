<template>
  <div class="page">
    <van-nav-bar title="我的借阅" left-arrow fixed placeholder @click-left="$router.back()" />
    <van-tabs v-model:active="activeRole" @change="onTabChange">
      <van-tab title="我借入的" name="borrower" />
      <van-tab title="我收到的申请" name="seller" />
    </van-tabs>

    <div v-if="list.length" class="borrow-list">
      <div v-for="r in list" :key="r.id" class="borrow-card" @click="goBook(r.book_id)">
        <div class="borrow-head">
          <span class="borrow-title van-ellipsis">{{ r.book?.title || '书籍' }}</span>
          <span class="borrow-tags">
            <van-tag v-if="r.overdue" type="danger" round>已逾期</van-tag>
            <van-tag :type="tagType(r.status)" round>{{ r.status_text }}</van-tag>
          </span>
        </div>
        <div class="borrow-meta">
          <span>{{ activeRole === 'borrower' ? `向 ${r.seller?.name || '卖家'} 借阅` : `借阅人：${r.borrower?.name || '同学'}` }}</span>
        </div>
        <div class="borrow-meta">
          <span v-if="r.due_at">到期日：<span :class="{ overdue: r.overdue }">{{ r.due_at }}</span></span>
          <span v-else>申请于 {{ r.created_at }}</span>
        </div>
        <div class="borrow-btns" @click.stop>
          <template v-if="activeRole === 'borrower'">
            <van-button v-if="r.status === 'approved'" size="small" type="primary" icon="revoke" :loading="loadingId === r.id" @click="onReturn(r)">归还</van-button>
            <van-button size="small" plain type="primary" icon="chat-o" @click="onContact(r)">联系{{ counterpartLabel }}</van-button>
          </template>
          <template v-else>
            <van-button v-if="r.status === 'pending'" size="small" type="success" :loading="loadingId === r.id" @click="onApprove(r)">同意</van-button>
            <van-button v-if="r.status === 'pending'" size="small" plain type="default" :loading="loadingId === r.id" @click="onReject(r)">拒绝</van-button>
            <van-button v-if="r.status === 'returned'" size="small" type="success" :loading="loadingId === r.id" @click="onConfirm(r)">确认拿回</van-button>
            <van-button v-if="r.status === 'approved' || r.status === 'returned'" size="small" plain type="warning" icon="bullhorn-o" :loading="loadingId === r.id" @click="onRemind(r)">发送提醒</van-button>
            <van-button size="small" plain type="primary" icon="chat-o" @click="onContact(r)">联系{{ counterpartLabel }}</van-button>
          </template>
        </div>
      </div>
    </div>
    <EmptyState v-else :description="activeRole === 'borrower' ? '暂无借阅记录' : '暂无收到的申请'" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { showConfirmDialog, showToast } from 'vant'
import { listBorrows, approveBorrow, rejectBorrow, returnBorrow, confirmReturnBorrow, remindBorrow } from '@/api/borrow'
import { createConversationFromBook, createConversationToUser } from '@/api/conversation'
import type { BorrowRequest } from '@/types'
import EmptyState from '@/components/EmptyState.vue'

const router = useRouter()
const activeRole = ref<'borrower' | 'seller'>('borrower')
const list = ref<BorrowRequest[]>([])
const loadingId = ref(0)
const counterpartLabel = computed(() => (activeRole.value === 'borrower' ? '卖家' : '借阅人'))

onMounted(loadList)

async function loadList() {
  try {
    const data = await listBorrows(activeRole.value, { page: 1, page_size: 50 })
    list.value = data.list
  } catch {
    list.value = []
  }
}

function onTabChange() {
  loadList()
}

function tagType(status: string): 'primary' | 'warning' | 'success' | 'default' | 'danger' {
  switch (status) {
    case 'approved':
      return 'primary'
    case 'pending':
      return 'warning'
    case 'completed':
      return 'success'
    case 'rejected':
      return 'default'
    case 'returned':
      return 'danger'
    default:
      return 'default'
  }
}

async function run(id: number, fn: () => Promise<unknown>, message: string) {
  try {
    loadingId.value = id
    await fn()
    showToast(message)
    await loadList()
  } catch {
    /* cancelled or handled */
  } finally {
    loadingId.value = 0
  }
}

async function onApprove(r: BorrowRequest) {
  try {
    await showConfirmDialog({ title: '同意借阅', message: `同意 ${r.borrower?.name || '该同学'} 借阅？同意后其他申请将自动拒绝。` })
  } catch {
    return
  }
  await run(r.id, () => approveBorrow(r.id), '已同意，书籍标记为借出')
}

async function onReject(r: BorrowRequest) {
  try {
    await showConfirmDialog({ title: '拒绝借阅', message: `拒绝 ${r.borrower?.name || '该同学'} 的借阅申请？` })
  } catch {
    return
  }
  await run(r.id, () => rejectBorrow(r.id), '已拒绝')
}

async function onReturn(r: BorrowRequest) {
  try {
    await showConfirmDialog({ title: '申请归还', message: '确认已将书归还给卖家？' })
  } catch {
    return
  }
  await run(r.id, () => returnBorrow(r.id), '已通知卖家，请等待确认')
}

async function onConfirm(r: BorrowRequest) {
  try {
    await showConfirmDialog({ title: '确认拿回', message: '确认已拿回这本书？确认后书籍恢复可借。' })
  } catch {
    return
  }
  await run(r.id, () => confirmReturnBorrow(r.id), '已确认拿回，书籍恢复可借')
}

async function onRemind(r: BorrowRequest) {
  try {
    await showConfirmDialog({ title: '发送提醒', message: r.overdue ? '书籍已逾期，确认发送还书提醒？' : '确认发送还书提醒？' })
  } catch {
    return
  }
  await run(r.id, () => remindBorrow(r.id), '提醒已通过站内消息发送')
}

async function onContact(r: BorrowRequest) {
  try {
    const conv =
      activeRole.value === 'borrower'
        ? await createConversationFromBook(r.book_id, '你好，关于借阅的书想和你沟通。')
        : await createConversationToUser(r.book_id, r.borrower_id, '你好，关于你借阅的书想和你沟通。')
    router.push(`/chat/${conv.id}`)
  } catch {
    /* toast handled */
  }
}

function goBook(bookId: number) {
  router.push(`/books/${bookId}`)
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
.borrow-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}
.borrow-title {
  font-size: 15px;
  font-weight: 600;
  flex: 1;
  min-width: 0;
}
.borrow-tags {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}
.borrow-meta {
  margin-top: 6px;
  font-size: 13px;
  color: #646566;
}
.borrow-meta .overdue {
  color: #ee0a24;
  font-weight: 600;
}
.borrow-btns {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  margin-top: 10px;
  flex-wrap: wrap;
}
</style>
