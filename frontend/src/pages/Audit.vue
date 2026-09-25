<template>
  <div class="page">
    <van-nav-bar title="操作审计日志" left-arrow fixed placeholder @click-left="$router.back()" />
    <van-field v-model="keyword" placeholder="按操作(action)筛选" clearable @keyup.enter="onFilter" />
    <van-list v-model:loading="loading" :finished="finished" finished-text="没有更多了" @load="onLoad">
      <van-cell-group inset v-for="a in list" :key="a.id" class="audit-card">
        <van-cell :title="a.action" :value="a.created_at" />
        <van-cell title="操作者" :value="a.user_id ? `用户 #${a.user_id}` : '未登录'" />
        <van-cell title="资源" :value="`${a.resource_type} / ${a.resource_id || '-'}`" />
        <van-cell title="请求ID" :label="a.request_id" />
        <van-cell title="IP" :value="a.ip" />
      </van-cell-group>
    </van-list>
    <EmptyState v-if="finished && !list.length" description="暂无审计记录" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { listAuditLogs } from '@/api/audit'
import type { AuditLog } from '@/types'
import EmptyState from '@/components/EmptyState.vue'

const list = ref<AuditLog[]>([])
const keyword = ref('')
const page = ref(1)
const total = ref(0)
const loading = ref(false)
const finished = ref(false)

async function fetchPage() {
  loading.value = true
  try {
    const data = await listAuditLogs({ page: page.value, page_size: 10, action: keyword.value || undefined })
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

fetchPage()
</script>

<style scoped>
.audit-card {
  margin-bottom: 8px;
}
</style>
