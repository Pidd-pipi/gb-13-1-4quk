<template>
  <div class="page">
    <van-nav-bar title="求购信息" fixed placeholder />
    <div class="filter-bar">
      <van-dropdown-menu>
        <van-dropdown-item v-model="query.subject_category" :options="subjectOptions" @change="onFilter" />
        <van-dropdown-item v-model="query.status" :options="statusOptions" @change="onFilter" />
      </van-dropdown-menu>
    </div>

    <van-list v-model:loading="loading" :finished="finished" finished-text="没有更多了" @load="onLoad">
      <van-cell-group v-for="w in list" :key="w.id" inset class="wish-card">
        <van-cell>
          <template #title>
            <div class="wish-title">{{ w.book_title }}</div>
            <div class="wish-meta">{{ w.author || '佚名' }} · {{ w.subject_text }} · {{ w.condition_text || '新旧不限' }}</div>
          </template>
          <template #value><span class="price">{{ formatPrice(w.expected_price) }}</span></template>
        </van-cell>
        <van-cell v-if="w.description" :title="w.description" />
        <van-cell :title="`发布人：${w.user?.name || ''}（${w.user?.campus || ''}）`" />
        <div class="wish-actions">
          <StatusBadge :status="w.status" kind="wish" />
          <van-button v-if="isMine(w) && w.status === 'open'" size="mini" @click="onClose(w)">关闭</van-button>
          <van-button v-if="!isMine(w) && w.status === 'open'" size="mini" type="primary" @click="onContact(w)">我有书，联系TA</van-button>
        </div>
      </van-cell-group>
    </van-list>
    <EmptyState v-if="finished && !list.length" description="暂无求购信息" />

    <van-button round block type="primary" class="fab" icon="plus" @click="$router.push('/publish-wish')">发布求购</van-button>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { showConfirmDialog, showToast } from 'vant'
import { listWishes, closeWish, contactWish } from '@/api/wish'
import type { Wish } from '@/types'
import { formatPrice } from '@/utils/format'
import { SubjectCategoryOptions } from '@/constants/enums'
import { useAuthStore } from '@/stores/authStore'
import StatusBadge from '@/components/StatusBadge.vue'
import EmptyState from '@/components/EmptyState.vue'

const router = useRouter()
const auth = useAuthStore()
const query = reactive({ subject_category: '', status: 'open' })
const list = ref<Wish[]>([])
const page = ref(1)
const total = ref(0)
const loading = ref(false)
const finished = ref(false)

const subjectOptions = [{ text: '全部分类', value: '' }, ...SubjectCategoryOptions.map((o) => ({ text: o.label, value: o.value }))]
const statusOptions = [
  { text: '进行中', value: 'open' },
  { text: '已关闭', value: 'closed' },
  { text: '全部', value: '' },
]

function isMine(w: Wish) {
  return auth.user?.id === w.user_id
}

async function fetchPage() {
  loading.value = true
  try {
    const data = await listWishes({
      subject_category: query.subject_category,
      status: query.status,
      page: page.value,
      page_size: 10,
    })
    total.value = data.total
    if (page.value === 1) {
      list.value = data.list
    } else {
      list.value = list.value.concat(data.list)
    }
    finished.value = list.value.length >= total.value
  } catch {
    finished.value = true
  } finally {
    loading.value = false
  }
}

function onLoad() {
  if (list.value.length >= total.value) return
  page.value += 1
  fetchPage()
}

function onFilter() {
  page.value = 1
  finished.value = false
  list.value = []
  fetchPage()
}

async function onClose(w: Wish) {
  try {
    await showConfirmDialog({ title: '关闭求购', message: '确认关闭这条求购信息？' })
    await closeWish(w.id)
    showToast('已关闭')
    onFilter()
  } catch {
    /* cancelled */
  }
}

async function onContact(w: Wish) {
  if (!auth.isLoggedIn) {
    router.push('/login')
    return
  }
  try {
    const conv = await contactWish(w.id, '你好，我有这本书，可以联系一下吗？')
    router.push(`/chat/${conv.id}`)
  } catch {
    /* toast handled */
  }
}

fetchPage()
</script>

<style scoped>
.wish-card {
  margin-bottom: 10px;
}
.wish-title {
  font-size: 15px;
  font-weight: 600;
}
.wish-meta {
  font-size: 12px;
  color: #969799;
  margin-top: 4px;
}
.price {
  color: #ee0a24;
  font-weight: 600;
}
.wish-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  padding: 8px 12px;
}
.fab {
  position: fixed;
  right: 16px;
  bottom: 76px;
  width: 120px;
  z-index: 10;
}
</style>
